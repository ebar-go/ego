package component

import "github.com/spf13/viper"

type Config struct {
	Named
	*viper.Viper
}

func NewConfig() *Config {
	return &Config{Named: componentConfig, Viper: viper.New()}
}
