package main

import (
	"context"
	"fmt"
	"github.com/wyy-go/wzo/carry"
	"github.com/wyy-go/wzo/core/transport/http/server/middleware/error_response"

	"github.com/gin-gonic/gin"
	"github.com/wyy-go/wzo"
	"github.com/wyy-go/wzo/core/log"
	"github.com/wyy-go/wzo/core/transport/http/server"
	"github.com/wyy-go/wzo/examples/proto/errapi"
	"github.com/wyy-go/wzo/examples/proto/errno"
)

// curl -w " status=%{http_code}" http://localhost:5180/error/internal
// curl -w " status=%{http_code}" http://localhost:5180/error/bad
// curl -w " status=%{http_code}" http://localhost:5180/error/biz
// curl -w " status=%{http_code}" http://localhost:5180/error/wzo

func main() {
	app := wzo.New(wzo.InitHttpServer(InitHttpServer))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type translatorData struct{}

func (t translatorData) TranslateData(v any) any {
	return &Result{
		Code:    200,
		Message: "ok",
		Data:    v,
	}
}

func InitHttpServer(s *server.Server) error {
	gin.DisableBindValidation()
	carrier := carry.NewCarry(carry.WithTranslatorData(translatorData{}))
	s.Use(server.CarrierInterceptor(carrier))
	s.Use(server.TransportInterceptor())
	s.Use(error_response.ErrorResponse())
	s.Use(func(c *gin.Context) {
		defer func() {
			v, ok := server.GetMetadata(c)
			log.Infof("---> %v %v", v, ok)
		}()
		c.Next()
	})

	g := s.Group("/")
	errapi.RegisterGreeterHTTPServer(g, &GreeterImpl{})

	return nil
}

type GreeterImpl struct{}

func (s *GreeterImpl) SayHello(ctx context.Context, req *errapi.HelloRequest) (*errapi.HelloReply, error) {
	rsp := &errapi.HelloReply{
		Message: fmt.Sprintf("hello %s!", req.Name),
	}

	return rsp, nil
}

func (s *GreeterImpl) TestError(ctx context.Context, req *errapi.ErrorRequest) (*errapi.ErrorReply, error) {
	if req.Name == "internal" {
		return nil, errno.ErrInternalServerw(errno.WithDetail("服务器错误详请"))
	} else if req.Name == "bad" {
		return nil, errno.ErrInternalServerw(errno.WithDetail("请求参数错误详请"))
	} else if req.Name == "biz" {
		return nil, errno.ErrBizError()
	}

	rsp := &errapi.ErrorReply{
		Message: fmt.Sprintf("[%s] 不是错误", req.Name),
	}
	return rsp, nil
}
