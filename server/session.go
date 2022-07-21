package server

import (
	"fmt"
	"passport-v4/config"
	"strconv"

	"github.com/gin-contrib/sessions"
	ss_redis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func initSession(e *gin.Engine) error {

	cfg := config.Redis
	uri := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	redis_store, err := ss_redis.NewStoreWithDB(1000, "tcp", uri, cfg.Password, strconv.Itoa(cfg.DB), config.Server.SessionSecret)
	if err != nil {
		logrus.Fatal("cannot connect to Redis! ", err.Error())
	}
	redis_store.Options(sessions.Options{
		Path:     "/",
		Domain:   config.Server.Domain,
		MaxAge:   config.Server.SessionExpire,
		Secure:   true,
		HttpOnly: true,
	})
	if err != nil {
		logrus.Error("cannot init redistore for gin session")
		return err
	}
	e.Use(sessions.Sessions("SESSION_TOKEN", redis_store))

	return nil
}
