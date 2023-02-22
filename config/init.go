package config

import (
	"encoding/base64"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Server   ServerType
	Redis    RedisType
	Mongo    MongoType
	ZjuOauth ZjuOauthType
)

func Init() {
	log.Info("[config] Init...")
	config := ConfigType{}
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error in ReadInConfig: %s", err)
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Error in Unmarshal: %s", err)
		return
	}

	Server = config.Server
	Redis = config.Redis
	Mongo = config.Mongo
	ZjuOauth = config.ZjuOauth

	Server.SessionSecret = make([]byte, 1000)
	_, err = base64.RawStdEncoding.Decode(Server.SessionSecret, []byte(Server.SessionSecretString))
	if err != nil {
		log.Warnf("SessionSecret set failed: %s", err)
	}
	log.Info("[Config] Init success")
}
