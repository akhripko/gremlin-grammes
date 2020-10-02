package options

import (
	"github.com/spf13/viper"
)

func ReadEnv() *Config {
	viper.AutomaticEnv()

	viper.SetEnvPrefix("APP")

	viper.SetDefault("GREMLIN_ADDR", "ws://127.0.0.1:8182")

	return &Config{
		GremlinAddr: viper.GetString("GREMLIN_ADDR"),
	}
}
