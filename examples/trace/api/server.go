package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wyy-go/wzo"
	"github.com/wyy-go/wzo/core/log"
	"github.com/wyy-go/wzo/core/transport/http/server"
	"github.com/wyy-go/wzo/core/transport/rpc/client"
	"github.com/wyy-go/wzo/examples/proto/greeter"
)

// curl http://127.0.0.1:5180/hello/wzo
func main() {
	app := wzo.New(wzo.InitHttpServer(InitHttpServer))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func InitHttpServer(r *server.Server) error {
	r.GET("/hello/:name", func(c *gin.Context) {
		cc, err := client.NewClient(
			client.WithServiceName("Greeter"),
			client.WithServiceAddr("127.0.0.1:5188"),
			client.Tracing(true),
		)
		if err != nil {
			log.Error(err)
			return
		}
		cli := greeter.NewGreeterClient(cc.GetXClient())

		args := &greeter.HelloRequest{
			Name: c.Param("name"),
		}

		log.Infof(args.Name)

		reply, err := cli.SayHello(c.Request.Context(), args)
		if err != nil {
			log.Error(err)
			return
		}

		c.String(http.StatusOK, reply.Message)
	})

	return nil
}
