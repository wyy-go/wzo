package main

import (
	"github.com/wyy-go/wzo"
	"github.com/wyy-go/wzo/core/log"
	"github.com/wyy-go/wzo/core/transport/http/server"
	"github.com/wyy-go/wzo/examples/layout/internal/app"
	"github.com/wyy-go/wzo/examples/layout/internal/router"
)

func main() {
	app := wzo.New(
		wzo.BeforeStart(app.Init),
		wzo.InitHttpServer(func(s *server.Server) error {
			router.Setup(s)
			return nil
		}))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
