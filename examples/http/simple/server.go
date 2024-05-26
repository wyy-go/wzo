package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/wyy-go/wzo"
	"github.com/wyy-go/wzo/carry"
	"github.com/wyy-go/wzo/core/log"
	"github.com/wyy-go/wzo/core/transport/http/server"
	"github.com/wyy-go/wzo/core/transport/http/server/middleware/error_response"
)

// curl http://127.0.0.1:5180/hello/wzo
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

	s.GET("/hello/:name", func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprintf("hello %s!", c.Param("name")))
	})
	return nil
}
