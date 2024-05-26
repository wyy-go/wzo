package main

import (
	"fmt"
	"github.com/wyy-go/wzo"
	"github.com/wyy-go/wzo/core/errors"
	"github.com/wyy-go/wzo/core/log"
	"github.com/wyy-go/wzo/core/transport/http/server"
)

func main() {
	app := wzo.New(wzo.InitHttpServer(InitHttpServer))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

type HelloRequest struct {
	Name string `json:"name,omitempty"`
}

type HelloResponse struct {
	Message string `json:"message,omitempty"`
}

func InitHttpServer(s *server.Server) error {
	s.PostEx("/hello", func(c *server.Context) {
		req := HelloRequest{}
		if err := c.ShouldBind(&req); err != nil {
			e := errors.FromError(err)
			c.JSON(500, e)
			c.Abort()
			return
		}

		rsp := HelloResponse{
			Message: fmt.Sprintf("hello %s!", req.Name),
		}
		c.JSON(200, rsp)
	})
	return nil
}
