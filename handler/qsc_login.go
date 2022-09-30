package handler

import (
	"net/http"
	"passport-v4/database"
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
		Username string `form:"username"`
		PassWord string `form:"password"`
	}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_WRONG_REQUEST, "参数错误")
		return
	}
	us, err := database.FindByQscId(model.User{Name: req.Username}) //这里的username实际上是qscid
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_DATABASE_ERROR, "数据库查找失败")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(us.PassWord), []byte(req.PassWord))
	user := model.QscProfile2User(us)
	if err != nil { //密码不匹配,返回一个空的user
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_WRONG_PASSWORD, "密码错误")
		return
	}
	ss := sessions.Default(c)
	ss.Set(SS_KEY_USER, user)
	ss.Save()
	resp.Json(c, user)
}

func SetPassword(c *gin.Context) {
	ss := sessions.Default(c)
	user := ss.Get(SS_KEY_USER).(model.User)
	if user.LoginType != model.LT_QSC {
		logrus.Warn("LoginType Error")
		c.AbortWithStatus(http.StatusBadRequest)
	}
	NewPassWord := c.Query("Password")
	qscuser, err := database.FindByName(user)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_DATABASE_ERROR, "数据库查找失败")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(NewPassWord), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_DBECRIPT_ERROR, "数据库加密失败")
		return
	}
	qscuser.PassWord = string(hash)
	err = database.UpdateUser(qscuser)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_DATABASE_ERROR, "数据库更新失败")
		return
	}
	resp.Json(c, nil)
}

func QscRegister(c *gin.Context) {
	var qscuser model.UserProfileQsc
	err := c.ShouldBindQuery(&qscuser)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_WRONG_REQUEST, "参数错误")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(qscuser.PassWord), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_DBECRIPT_ERROR, "数据库加密失败")
		return
	}
	qscuser.PassWord = string(hash)
	err = database.Insert(qscuser)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_DATABASE_ERROR, "数据库插入失败")
		return
	}
	resp.Json(c, qscuser)
}
