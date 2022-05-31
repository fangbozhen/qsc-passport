package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Init(e *gin.Engine) error {
	logrus.Info("[server] Init...")
	if err := initSession(e); err != nil {
		return err
	}
	configRoutes(e)
	logrus.Info("[server] Init Success")
	return nil
}
