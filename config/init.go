package config

import (
	"encoding/base64"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() error {
	logrus.Info("[config] Init...")
	config := ConfigType{}
	var err error = nil

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}

	Server = config.Server
	Redis = config.Redis
	ZjuOauth = config.ZjuOauth
	Mongo = config.Mongo

	Server.SessionSecret = make([]byte, 1000)
	_, err = base64.RawStdEncoding.Decode(Server.SessionSecret, []byte(Server.SessionSecretString))
	if err != nil {
		return err
	}
	logrus.Info("[config] Init Success")
	return nil
}
