package core

import (
	"errors"
	"fmt"
	"os"

	"github.com/slntopp/core-chatting/cc"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	CONFIG_LOCATION string
	ROOT_ADMIN      string
)

func init() {
	viper.AutomaticEnv()
	viper.SetDefault("CONFIG_LOCATION", "cc.yaml")
	viper.SetDefault("ROOT_ADMIN", "0")

	CONFIG_LOCATION = viper.GetString("CONFIG_LOCATION")
	ROOT_ADMIN = viper.GetString("ROOT_ADMIN")

	_, err := os.ReadFile(CONFIG_LOCATION)
	if err != nil {
		panic(fmt.Errorf("couldn't read config. Location: %s. Error: %v", CONFIG_LOCATION, err))
	}
}

func Config() (*cc.Defaults, error) {
	var config cc.Defaults

	conf, err := os.ReadFile(CONFIG_LOCATION)
	if err != nil {
		return &config, err
	}

	err = yaml.Unmarshal(conf, &config)
	return &config, err
}

func SetConfig(requestor string, defaults *cc.Defaults) (*cc.Defaults, error) {
	if requestor != ROOT_ADMIN {
		return nil, errors.New("not root account")
	}

	marshal, err := yaml.Marshal(defaults)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(CONFIG_LOCATION, marshal, 0666)
	if err != nil {
		return nil, err
	}

	return defaults, err
}
