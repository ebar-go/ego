package config

import (
	"github.com/spf13/viper"
)

type Instance struct {
	*viper.Viper
}

// LoadFile loads the configuration file specified.
func (c *Instance) LoadFile(paths ...string) error {
	for _, p := range paths {
		c.SetConfigFile(p)
		if err := c.MergeInConfig(); err != nil {
			return err
		}
	}
	return nil
}

func New() *Instance {
	return &Instance{Viper: viper.New()}
}
