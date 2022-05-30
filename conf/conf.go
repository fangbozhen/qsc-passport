package conf

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConf() {
	viper.SetConfigName("conf")
	viper.AddConfigPath("./")

	// TODO:

	logrus.Info("Finish initialization of config file.")
}
