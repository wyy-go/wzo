package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wyy-go/wzo/core/log"
	"github.com/wyy-go/wzo/core/pkg/addr"
	znet "github.com/wyy-go/wzo/core/pkg/net"
	"github.com/wyy-go/wzo/core/registry"
	"github.com/wyy-go/wzo/core/transport/http/server/middleware/logging"
	"github.com/wyy-go/wzo/core/transport/http/server/middleware/tracing"
)

type Server struct {
	opts Options
	*gin.Engine
	server *http.Server
}

func NewServer(opts ...Option) *Server {
	options := newOptions(opts...)

	srv := &Server{
		opts: options,
	}

	gin.SetMode(srv.opts.Mode)
	r := gin.New()
	srv.Engine = r
	srv.server = &http.Server{Handler: srv.Engine}
	return srv
}

func (s *Server) Init(opts ...Option) {
	for _, opt := range opts {
		opt(&s.opts)
	}
}

func (s *Server) Start() error {
	s.Engine.Use(TransportInterceptor())

	if s.opts.Tracing {
		s.Engine.Use(tracing.Trace(s.opts.Name))
	}

	s.Engine.Use(logging.Log())

	if s.opts.InitHttpServer != nil {
		if err := s.opts.InitHttpServer(s); err != nil {
			return err
		}
	}

	l, err := net.Listen("tcp", s.opts.Addr)
	if err != nil {
		return err
	}
	a := l.Addr().String()
	log.Infof("Server [GIN] listening on %s", a)

	s.register(a)

	go func() {
		if err := s.server.Serve(l); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err.Error())
			}
		}
	}()
	return nil
}

func (s *Server) register(a string) {
	if s.opts.Registry == nil {
		return
	}
	var err error
	var host, port string
	if cnt := strings.Count(a, ":"); cnt >= 1 {
		host, port, err = net.SplitHostPort(a)
		if err != nil {
			log.Fatal(err)
			return
		}
	} else {
		host = a
	}

	address, err := addr.Extract(host)
	if err != nil {
		log.Fatal(err)
	}

	if port != "" {
		address = znet.HostPort(address, port)
	}

	if err := s.opts.Registry.Register(&registry.Service{
		Scheme:  "http",
		Name:    s.opts.Name,
		Address: address,
		Weight:  0,
	}); err != nil {
		log.Fatal(err)
	}

	s.opts.Addr = address

	log.Infof("Registering server: %s", address)
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
