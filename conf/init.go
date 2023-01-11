package conf

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Server ServerType
	Redis  RedisType
	Mongo  MongoType
)

func Init() {
	log.Info("[config] Init...")
	config := ConfigType{}
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigName("yaml")

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

	log.Info("[Config] Init success")
}
