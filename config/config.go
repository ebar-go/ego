package config

import (
	"github.com/spf13/viper"
	"go.uber.org/dig"
)

var Container = dig.New()

// ReadFromEnvironment read from system environment
func ReadFromEnvironment() {
	viper.AutomaticEnv()
}

// ReadFromFile read from file
func ReadFromFile(path string) error {
	viper.SetConfigFile(path)

	return viper.ReadInConfig()
}
