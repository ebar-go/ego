package config

// MnsConfig 阿里云MNS 配置项
type MnsConfig struct {
	Url             string `json:"url"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}
