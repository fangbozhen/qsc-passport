package main

import (
	"QSCpassport/database"
	"QSCpassport/server"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"time"

	"QSCpassport/config"
	"QSCpassport/model"
)

func initLogger() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceQuote:      true,
		ForceColors:     true,
	})
}

func initSentry() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://19fc972b48574d11920c6b72e9c39af6@sentry.zjuqsc.com/5",
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}

func main() {
	initLogger()
	initSentry()
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
