package main

import (
	"context"
	"fmt"

	"github.com/smallnest/rpcx/server"
	"github.com/wyy-go/wzo"
	"github.com/wyy-go/wzo/core/log"
	"github.com/wyy-go/wzo/examples/proto/greeter"
)

func main() {
	app := wzo.New(wzo.InitRpcServer(InitRpcServer))

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
