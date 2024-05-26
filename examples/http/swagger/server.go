package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/wyy-go/wzo"
	"github.com/wyy-go/wzo/carry"
	"github.com/wyy-go/wzo/core/log"
	"github.com/wyy-go/wzo/core/transport/http/server"
	"github.com/wyy-go/wzo/core/transport/http/server/middleware/error_response"
	"github.com/wyy-go/wzo/examples/proto/greeter"
)

// curl http://127.0.0.1:5180/hello/wzo
// http://127.0.0.1:5180/swagger/index.html
func main() {
	app := wzo.New(wzo.InitHttpServer(InitHttpServer))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func InitHttpServer(s *server.Server) error {
	gin.DisableBindValidation()
	carrier := carry.NewCarry()
	s.Use(server.CarrierInterceptor(carrier))
	s.Use(error_response.ErrorResponse())

	Swagger(s)
	g := s.Group("/")
	greeter.RegisterGreeterHTTPServer(g, &GreeterImpl{})

	return nil
}

type GreeterImpl struct {
	greeter.GreeterHTTPServer
}

func (s *GreeterImpl) SayHello(ctx context.Context, req *greeter.HelloRequest) (*greeter.HelloReply, error) {
	rsp := &greeter.HelloReply{
		Message: fmt.Sprintf("hello %s!", req.Name),
	}

	return rsp, nil
}
