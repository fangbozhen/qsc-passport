package server

import (
	"fmt"
	"passport-v4/config"
	"strconv"

	"github.com/gin-contrib/sessions"
	ss_redis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func initSession(e *gin.Engine) error {

	cfg := config.Redis
	uri := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	redis_store, err := ss_redis.NewStoreWithDB(1000, "tcp", uri, cfg.Password, strconv.Itoa(cfg.DB), config.Server.SessionSecret)
	redis_store.Options(sessions.Options{
		Path:     "/",
		Domain:   config.Server.Domain,
		MaxAge:   7200,
		Secure:   true,
		HttpOnly: true,
	})
	if err != nil {
		return err
	}
	e.Use(sessions.Sessions("SESSION_TOKEN", redis_store))

	return nil
}
