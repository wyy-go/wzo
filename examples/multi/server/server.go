package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/smallnest/rpcx/server"

	"github.com/wyy-go/wzo"
	"github.com/wyy-go/wzo/carry"
	"github.com/wyy-go/wzo/core/log"
	httpserver "github.com/wyy-go/wzo/core/transport/http/server"
	"github.com/wyy-go/wzo/examples/proto/greeter"
)

func main() {
	app := wzo.New(
		wzo.InitRpcServer(InitRpcServer),
		wzo.InitHttpServer(InitHttpServer),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func InitRpcServer(s *server.Server) error {
	if err := s.RegisterName("Greeter", &GreeterImpl{}, ""); err != nil {
		return err
	}
	return nil
}

type GreeterImpl struct{}

func (s *GreeterImpl) SayHello(ctx context.Context, req *greeter.HelloRequest, rsp *greeter.HelloReply) error {
	*rsp = greeter.HelloReply{
		Message: fmt.Sprintf("hello %s!", req.Name),
	}

	return nil
}

func InitHttpServer(r *httpserver.Server) error {
	gin.DisableBindValidation()
	carrier := carry.NewCarry()
	r.Use(httpserver.CarrierInterceptor(carrier))

	g := r.Group("/")
	greeter.RegisterGreeterHTTPServer(g, &HttpGreeter{})

	return nil
}

type HttpGreeter struct{}

func (s *HttpGreeter) SayHello(ctx context.Context, req *greeter.HelloRequest) (*greeter.HelloReply, error) {
	rsp := &greeter.HelloReply{
		Message: fmt.Sprintf("hello %s!", req.Name),
	}

	return rsp, nil
}
