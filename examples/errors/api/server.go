package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wyy-go/wzo/carry"

	"github.com/wyy-go/wzo"
	"github.com/wyy-go/wzo/core/errors"
	"github.com/wyy-go/wzo/core/log"
	"github.com/wyy-go/wzo/core/transport/http/server"
	"github.com/wyy-go/wzo/core/transport/http/server/middleware/error_response"
	"github.com/wyy-go/wzo/core/transport/rpc/client"
	"github.com/wyy-go/wzo/examples/proto/errapi"
)

// curl -w " status=%{http_code}" http://localhost:5180/error/internal
// curl -w " status=%{http_code}" http://localhost:5180/error/bad
// curl -w " status=%{http_code}" http://localhost:5180/error/biz
// curl -w " status=%{http_code}" http://localhost:5180/error/wzo

var cc *client.Client

func main() {
	app := wzo.New(wzo.InitHttpServer(InitHttpServer))

	cc, _ = client.NewClient(client.WithServiceName("ErrAPI"), client.WithServiceAddr("127.0.0.1:5188"))

	if cc == nil {
		log.Fatal("err")
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func InitHttpServer(s *server.Server) error {
	gin.DisableBindValidation()
	carrier := carry.NewCarry()
	s.Use(server.CarrierInterceptor(carrier))
	s.Use(server.TransportInterceptor())
	s.Use(error_response.ErrorResponse())
	g := s.Group("")
	errapi.RegisterGreeterHTTPServer(g, &ErrImpl{})

	return nil
}

type ErrImpl struct{}

var _ errapi.GreeterHTTPServer = &ErrImpl{}

func (s *ErrImpl) SayHello(ctx context.Context, req *errapi.HelloRequest) (*errapi.HelloReply, error) {
	cli := errapi.NewGreeterClient(cc.GetXClient())

	request := errapi.HelloRequest{Name: req.Name}

	reply, err := cli.SayHello(ctx, &request)
	if err != nil {
		err = errors.WrapRpcError(err)
		log.Errorf("%+v", err)
		return nil, err
	}

	rsp := &errapi.HelloReply{
		Message: fmt.Sprintf("hello %s!", reply.Message),
	}
	return rsp, nil
}

func (s *ErrImpl) TestError(ctx context.Context, req *errapi.ErrorRequest) (*errapi.ErrorReply, error) {
	cli := errapi.NewGreeterClient(cc.GetXClient())
	request := errapi.ErrorRequest{Name: req.Name}

	reply, err := cli.TestError(ctx, &request)
	if err != nil {
		err = errors.WrapRpcError(err)
		log.Errorf("%+v", err)
		return nil, err
	}

	rsp := &errapi.ErrorReply{
		Message: reply.Message,
	}
	return rsp, nil
}
