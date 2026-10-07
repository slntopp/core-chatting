package core

import (
	"errors"
	"fmt"
	"os"
	"slices"

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

// ErrNotConfigAdmin is returned to anyone who may not rewrite the defaults.
var ErrNotConfigAdmin = errors.New("only the root account or a chatting admin may change the settings")

// CanSetConfig says whether requestor may rewrite the defaults: the root account,
// or one of the global admins already listed in them. A NoCloud admin who is not
// root was refused before, which left the Bot tab usable by one person only.
// Membership is read from the current file, never from the request, so nobody
// can list themselves in the same call that is being checked.
func CanSetConfig(requestor string) (bool, error) {
	if requestor == "" {
		return false, nil
	}
	if requestor == ROOT_ADMIN {
		return true, nil
	}
	current, err := Config()
	if err != nil {
		return false, err
	}
	return slices.Contains(current.GetAdmins(), requestor), nil
}

func SetConfig(requestor string, defaults *cc.Defaults) (*cc.Defaults, error) {
	allowed, err := CanSetConfig(requestor)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrNotConfigAdmin
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
