package runtime

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/wyy-go/wzo/core/config"
	zdb "github.com/wyy-go/wzo/examples/layout/pkg/database/db"
	zredis "github.com/wyy-go/wzo/examples/layout/pkg/database/redis"
	"gorm.io/gorm"
)

var configProviderSet = wire.NewSet(
	config.Default,
	NewConfig,
	wire.FieldsOf(
		new(*Config),
		"Mysql",
		"Redis",
	),
)

type Config struct {
	Mysql zdb.Config    `yaml:"mysql"`
	Redis zredis.Config `yaml:"redis"`
}

func NewConfig(c config.Config) (*Config, error) {
	var conf Config
	err := c.Unmarshal(&conf)
	return &conf, err
}

var RuntimeProviderSet = wire.NewSet(
	configProviderSet,
	zdb.Open,
	zredis.NewClient,
	wire.NewSet(wire.Struct(new(Runtime), "*")),
)

type Runtime struct {
	Config *Config
	DB     *gorm.DB
	RC     *redis.Client
}
