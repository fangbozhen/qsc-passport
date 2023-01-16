package server

import (
	"QSCpassport/config"
	"QSCpassport/server/handlers"
	"QSCpassport/server/middleware"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func configRoutes(r *gin.Engine) {

	r.Use(middleware.Response)
	r.GET("/ping", handlers.Ping)

	r.POST("/qsc/login", handlers.QscLoginJson)
	r.POST("/qsc/reset-password", handlers.SetPasswordJson)

	r.GET("/zju/login", handlers.ZjuOauthRequest)
	r.GET("/zju/login-success", handlers.ZjuOauthCodeReturn)

	r.GET("/logout", handlers.Logout)
	r.GET("/profile", handlers.GetProfile)

	//admin := r.Group("/admin")
	//{
	//	admin.GET("/login")
	//	admin.POST("/login")
	//	admin.GET("/index")
	//}
}

func initSession(r *gin.Engine) {
	log.Info("[Server] Session Init...")
	cfg := config.Redis
	uri := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	store, err := redis.NewStoreWithDB(1000, "tcp", uri, cfg.Password, strconv.Itoa(cfg.DB), config.Server.SessionSecret)
	if err != nil {
		log.Fatalf("Cannot connect to Redis: %s", err.Error())
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
	store.Options(opt)
	r.Use(sessions.Sessions("SessionToken", store))
}

func initCors(r *gin.Engine) {
	corsCfg := cors.Config{
		AllowOrigins:     []string{"https://www.qsc.zju.edu.cn"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           3600,
		ExposeHeaders:    []string{"Authorization", "Set-Cookie"},
	}
	r.Use(cors.New(corsCfg))
}

func Init(r *gin.Engine) {
	log.Info("[Server] Init...")
	initSession(r)
	initCors(r)
	configRoutes(r)
	log.Info("[Server] Init success")
}
