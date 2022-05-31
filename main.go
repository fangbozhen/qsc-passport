package main

import (
	"fmt"
	"math/rand"
	"os"
	"passport-v4/config"
	"passport-v4/database"
	"passport-v4/model"
	"passport-v4/server"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// err handler
func Do(err error) {
	if err != nil {
		logrus.Fatalf("Error in Init: %s", err)
	}
}

func initLogger() error {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:           "2006-01-02 15:04:05",
		FullTimestamp:             true,
		ForceQuote:                true,
		EnvironmentOverrideColors: false,
	})
	return nil
}

func main() {
	initLogger()
	Do(config.Init())
	Do(database.Init())
	Do(model.Init())

	e := gin.Default()
	Do(server.Init(e))
	rand.Seed(time.Hour.Milliseconds())

	logrus.Infof("Gin Server Started")
	err := e.Run(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port))
	if err != nil {
		logrus.Errorf("Error While Running Server: %s", err.Error())
		os.Exit(-1)
	}
}
