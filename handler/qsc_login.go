package handler

import (
	. "passport-v4/global"
	"passport-v4/model"
	"passport-v4/util/resp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func QscLogin(c *gin.Context) {
	var req struct {
		Username string
		Password string
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_WRONG_REQUEST, "参数错误")
		return
	}
	qscer, err := model.FindQSCerByQscId(req.Username)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_WRONG_USERNAME, "找不到用户名")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(qscer.Password), []byte(req.Password))
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_WRONG_PASSWORD, "密码错误")
		return
	}

	user := model.QscProfile2User(qscer)
	ss := sessions.Default(c)
	ss.Set(SS_KEY_USER, user)
	ss.Save()
	resp.Json(c, user)
}

func SetPassword(c *gin.Context) {
	ss := sessions.Default(c)
	user := ss.Get(SS_KEY_USER).(model.User)
	if user.LoginType != model.LT_QSC || user.QscUser == nil {
		logrus.Warn("LoginType Error")
		resp.Err(c, resp.E_AUTH_FAILED, "not qscer")
		return
	}
	var req struct {
		Password string
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_WRONG_REQUEST, "参数错误")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_DBECRIPT_ERROR, "bcrypt加密失败")
		return
	}
	user.QscUser.Password = string(hash)
	err = model.UpdateQSCer(*user.QscUser)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_DATABASE_ERROR, "更新密码失败")
		return
	}
	resp.Json(c, nil)
}

func QscRegister(c *gin.Context) {
	var qscuser model.UserProfileQsc
	err := c.ShouldBindJSON(&qscuser)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_WRONG_REQUEST, "参数错误")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(qscuser.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_DBECRIPT_ERROR, "加密失败")
		return
	}
	qscuser.Password = string(hash)
	err = model.InsertQSCer(qscuser)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_DATABASE_ERROR, "数据库插入失败")
		return
	}
	resp.Json(c, qscuser)
}
