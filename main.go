package main

import (
	"fmt"
	"math/rand"
	"os"
	"passport-v4/database"
	"passport-v4/server"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"passport-v4/config"
	"passport-v4/model"
)

func initLogger() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceQuote:      true,
		ForceColors:     true,
	})
}

func main() {
	initLogger()
	config.Init()
	model.Init()
	database.InitDb()
	router := gin.Default()
	server.Init(router)

	rand.Seed(time.Hour.Milliseconds())

	log.Info("Gin Server Started")
	err := router.Run(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port))
	if err != nil {
		log.Errorf("Error while running Server: %s", err.Error())
		os.Exit(-1)
	}
}
