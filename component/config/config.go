package config

import (
	"github.com/spf13/viper"
)

// Config 配置
type Config struct {
	*viper.Viper
}


// New 实例
func New() *Config {
	conf := new(Config)
	conf.Viper = viper.New()
	return conf
}


// LoadFile 加载配置文件
func (conf *Config) LoadFile(path ...string) error {
	for _, p := range path {
		conf.SetConfigFile(p)
		if err := conf.MergeInConfig(); err != nil {
			return err
		}
	}

	return nil
}


