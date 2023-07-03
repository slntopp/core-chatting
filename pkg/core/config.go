package core

import (
	"fmt"
	"os"

	"github.com/slntopp/core-chatting/cc"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var CONFIG_LOCATION string

func init() {
	viper.AutomaticEnv()
	viper.SetDefault("CONFIG_LOCATION", "cc.yaml")

	CONFIG_LOCATION = viper.GetString("CONFIG_LOCATION")

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
