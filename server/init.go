package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Init(e *gin.Engine) error {
	logrus.Info("[server] Init...")
	if err := initSession(e); err != nil {
		return err
	}
	cors_cfg := cors.Config{
		AllowOrigins:     []string{"https://www.qsc.zju.edu.cn"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           3600,
		ExposeHeaders:    []string{"Authorization", "Set-Cookie"},
	}
	e.Use(cors.New(cors_cfg))
	configRoutes(e)
	logrus.Info("[server] Init Success")
	return nil
}
