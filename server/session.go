package server

import (
	"fmt"
	"net/http"
	"passport-v4/config"
	"strconv"

	"github.com/gin-contrib/sessions"
	SessionStoreRedis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func initSession(e *gin.Engine) error {

	cfg := config.Redis
	uri := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	redisStore, err := SessionStoreRedis.NewStoreWithDB(1000, "tcp", uri, cfg.Password, strconv.Itoa(cfg.DB), config.Server.SessionSecret)
	if err != nil {
		logrus.Fatal("cannot connect to Redis! ", err.Error())
	}
	opt := sessions.Options{
		Path:     "/",
		Domain:   config.Server.Domain,
		MaxAge:   config.Server.SessionExpire,
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	}
	if gin.Mode() != gin.ReleaseMode {
		opt.Secure = false
	}
	redisStore.Options(opt)
	if err != nil {
		logrus.Error("cannot init redistore for gin session")
		return err
	}
	e.Use(sessions.Sessions("SESSION_TOKEN", redisStore))

	return nil
}
