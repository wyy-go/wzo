package app

import (
	"github.com/google/wire"
	"github.com/wyy-go/wzo/examples/layout/internal/service/example"
)

var serviceProviderSet = wire.NewSet(
	example.New,
	wire.NewSet(wire.Struct(new(Service), "*")),
)

type Service struct {
	*example.Example
}
