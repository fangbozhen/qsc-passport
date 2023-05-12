package server

import (
	"fmt"
	"net/http"
	"passport-v4/config"
	"passport-v4/server/handlers"
	"passport-v4/server/middleware"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func configRoutes(r *gin.Engine) {

	r.Use(middleware.Response)
	r.GET("/ping", handlers.Ping)

	r.GET("/qsc/login", handlers.QscLoginRediect)
	r.POST("/qsc/login", handlers.QscLoginJson)
	r.GET("/qsc/reset_password", handlers.SetPasswordRediect)
	r.POST("/qsc/reset_password", handlers.SetPasswordJson)

	r.GET("/zju/login", handlers.ZjuOauthRequest)
	r.GET("/zju/login_success", handlers.ZjuOauthCodeReturn)

	r.GET("/logout", handlers.Logout)
	r.GET("/profile", handlers.GetProfile)

	r.GET("/privilege", handlers.GetPivilege)

	admin := r.Group("/admin", middleware.AdminCheck)
	{
		admin.POST("/user/register", handlers.QscRegister)
		admin.POST("/user/upload", handlers.Upload)
		admin.POST("/user/updateone", handlers.UpdateOne)
		admin.POST("/user/updatemany", handlers.UpdateMany)
		admin.POST("/user/delete", handlers.Delete)
		admin.POST("/user/list", handlers.GetByPage)
	}
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
		SameSite: http.SameSiteNoneMode,
	}

	if gin.Mode() != gin.ReleaseMode {
		opt.Secure = false
	}
	store.Options(opt)
	r.Use(sessions.Sessions("SESSION_TOKEN", store))
}

func initCors(r *gin.Engine) {
	corsCfg := cors.Config{
		AllowOrigins:     []string{"https://www.qsc.zju.edu.cn", "https://www.zjuqsc.com"},
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
