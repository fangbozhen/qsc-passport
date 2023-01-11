package main

import (
	"QSCpassport/server"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"time"

	"QSCpassport/conf"
	"QSCpassport/models"
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
	conf.Init()
	models.Init()
	router := gin.Default()
	server.Init(router)

	rand.Seed(time.Hour.Milliseconds())

	log.Info("Gin Server Started")
	err := router.Run(fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port))
	if err != nil {
		log.Errorf("Error while running Server: %s", err.Error())
		os.Exit(-1)
	}
}
