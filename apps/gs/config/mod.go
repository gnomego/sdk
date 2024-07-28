package config

import (
	config3 "github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

var (
	Config   = config3.New("gs")
	dbConfig DbConfig
)

type DbConfig struct {
	Driver     string
	Connection string
}

func init() {
	Config.WithOptions(config3.ParseEnv)
	Config.AddDriver(yaml.Driver)
	Config.SetData(map[string]interface{}{
		"db": map[string]interface{}{
			"driver":     "sqlite",
			"connection": "file::memory:?cache=shared",
		},
	})
}
