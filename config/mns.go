package config

import "github.com/spf13/viper"

// MnsOptions 阿里云MNS 配置项
type mns struct {
	Url             string
	AccessKeyId     string
	AccessKeySecret string
}

const (
	mnsEndpointKey = "mns.endpoint"
	mnsAccessIdKey = "mns.accessId"
	mnsSecretKey = "mns.secret"
)

// Mns
func Mns() (options *mns) {
	if err := Container.Invoke(func(o *mns) {
		options = o
	}); err != nil {
		options = &mns{
			Url:             viper.GetString(mnsEndpointKey),
			AccessKeyId:     viper.GetString(mnsAccessIdKey),
			AccessKeySecret: viper.GetString(mnsSecretKey),
		}

		_ = Container.Provide(func() *mns {
			return options
		})
	}
	return
}