package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestLoadYaml
func TestLoadYaml(t *testing.T) {
	var conf = struct {
		Apollo struct{
			AppId string `yaml:"appId"`
			Cluster string
			Namespace string
			Ip string
			BackupConfigPath string `yaml:"backConfigPath"`
		}
	}{}

	err := LoadYaml(&conf, "/tmp/app.yaml")
	assert.Nil(t, err)
	fmt.Println(conf)
}

// TestLoadIni
func TestLoadIni(t *testing.T) {
	conf := struct {
		DB struct{
			Host string
			Port int
		}
	}{}

	err := LoadIni(&conf, "/tmp/app.ini")
	assert.Nil(t, err)
}

// TestLoadJson
func TestLoadJson(t *testing.T) {
	conf := struct {
		DB struct{
			Host string `json:"host"`
			Port int `json:"port"`
		} `json:"db"`
	}{}

	err := LoadJson(&conf, "/tmp/app.json")
	assert.Nil(t, err)
}
