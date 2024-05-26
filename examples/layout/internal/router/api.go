package router

import (
	"github.com/wyy-go/wzo/core/transport/http/server"
	"github.com/wyy-go/wzo/examples/layout/internal/app"
	"github.com/wyy-go/wzo/examples/proto/example"
)

func RegisterAPI(s *server.Server) {
	appContext := app.Context()
	g := s.Group("")
	example.RegisterExampleHTTPServer(g, appContext.Service.Example)
}
