package handlers

import (
	"fmt"
	"passport-v4/model"
	"passport-v4/utils"
	"passport-v4/utils/resp"

	"github.com/getsentry/sentry-go"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	BaseURL         = "https://www.qsc.zju.edu.cn/passport/v4/static/index.html#"
	LoginPath       = "/login"
	SetPasswordPath = "/change_password"
)

func QscLoginJson(c *gin.Context) {
	var req struct {
		QscId    string `json:"qscid"`
		Password string `json:"password"`
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Errorf("Request Error: %s", err.Error())
		resp.Err(c, resp.WrongRequestError, "参数错误")
		return
	}

	qscer, err := model.FindQSCerByQscId(req.QscId)
	if err != nil {
		log.Errorf("Request Error: %s", err.Error())
		resp.Err(c, resp.WrongUsernameError, "找不到用户名")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(qscer.Password), []byte(req.Password))
	if err != nil {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetUser(sentry.User{ID: qscer.QscId})
			sentry.CaptureMessage(fmt.Sprintf("qsc-login-password-error"))
		})
		log.Errorf("Request Error: %s", err.Error())
		resp.Err(c, resp.WrongPasswordError, "密码错误")
		return
	}

	user := model.QscProfile2User(qscer)
	ss := sessions.Default(c)
	ss.Set(utils.SessionKeyUser, user)
	err = ss.Save()
	if err != nil {
		log.Error(err)
	}
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{ID: qscer.QscId})
		sentry.CaptureMessage(fmt.Sprintf("qsc-login-success"))
	})
	resp.Json(c, user)
}

func SetPasswordJson(c *gin.Context) {
	ss := sessions.Default(c)
	if ss.Get(utils.SessionKeyUser) == nil {
		resp.Err(c, resp.AuthFailedError, "未登录！")
		return
	}
	user := ss.Get(utils.SessionKeyUser).(model.User)
	if user.LoginType != model.LoginQSC || user.QscUser == nil {
		log.Warn("LoginType Error")
		resp.Err(c, resp.AuthFailedError, "not qscer")
		return
	}
	var req struct {
		Password string
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Errorf("err: %s", err.Error())
		resp.Err(c, resp.WrongRequestError, "参数错误")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("err: %s", err.Error())
		resp.Err(c, resp.InternalError, "bcrypt加密失败")
		return
	}
	user.QscUser.Password = string(hash)
	err = model.UpdateQSCer(*user.QscUser)
	if err != nil {
		log.Errorf("err: %s", err.Error())
		resp.Err(c, resp.DatabaseError, "更新密码失败")
		return
	}
	resp.Json(c, nil)
}

func QscRegister(c *gin.Context) {
	var qscuser model.UserProfileQsc
	err := c.ShouldBindJSON(&qscuser)
	if err != nil {
		log.Errorf("err: %s", err.Error())
		resp.Err(c, resp.WrongRequestError, "参数错误")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(qscuser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("err: %s", err.Error())
		resp.Err(c, resp.InternalError, "加密失败")
		return
	}
	qscuser.Password = string(hash)
	err = model.InsertQSCer(qscuser)
	if err != nil {
		log.Errorf("err: %s", err.Error())
		resp.Err(c, resp.DatabaseError, "数据库插入失败")
		return
	}
	resp.Json(c, qscuser)
}

func QscLoginRediect(c *gin.Context) {
	c.Redirect(302, BaseURL+LoginPath+"?"+c.Request.URL.RawQuery)
}

func SetPasswordRediect(c *gin.Context) {
	c.Redirect(302, BaseURL+SetPasswordPath+"?"+c.Request.URL.RawQuery)
}
