package main

import (
	"context"
	"fmt"

	"github.com/smallnest/rpcx/server"

	"github.com/wyy-go/wzo"
	"github.com/wyy-go/wzo/core/log"
	"github.com/wyy-go/wzo/examples/proto/errapi"
	"github.com/wyy-go/wzo/examples/proto/errno"
)

func main() {
	app := wzo.New(wzo.InitRpcServer(InitRpcServer))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func InitRpcServer(s *server.Server) error {
	if err := s.RegisterName("ErrAPI", &GreeterImpl{}, ""); err != nil {
		return err
	}
	return nil
}

type GreeterImpl struct{}

var _ errapi.GreeterAble = &GreeterImpl{}

func (s *GreeterImpl) SayHello(ctx context.Context, req *errapi.HelloRequest, rsp *errapi.HelloReply) error {
	*rsp = errapi.HelloReply{
		Message: fmt.Sprintf("hello %s!", req.Name),
	}

	return nil
}

func (s *GreeterImpl) TestError(ctx context.Context, req *errapi.ErrorRequest, rsp *errapi.ErrorReply) error {
	if req.Name == "internal" {
		return errno.ErrInternalServerw(errno.WithDetail("服务器错误详请"))
	} else if req.Name == "bad" {
		return errno.ErrBadRequestw(errno.WithDetail("请求参数错误详请"))
	} else if req.Name == "biz" {
		return errno.ErrBizError()
	}

	*rsp = errapi.ErrorReply{
		Message: fmt.Sprintf("[%s] 不是错误", req.Name),
	}

	return nil
}
