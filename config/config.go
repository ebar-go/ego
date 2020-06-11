package config

import (
	"github.com/spf13/viper"
	"go.uber.org/dig"
)

var container = dig.New()

func init() {
	viper.AutomaticEnv()
}

// ReadFromFile read from file
func ReadFromFile(path string) error {
	viper.SetConfigFile(path)

	return viper.ReadInConfig()
}
