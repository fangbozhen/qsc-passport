package handlers

import (
	"QSCpassport/models"
	"QSCpassport/utils"
	"QSCpassport/utils/resp"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func QscLoginJson(c *gin.Context) {
	var req models.UserReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Errorf("Request Error: %s", err.Error())
		resp.Err(c, resp.WrongRequestError, "参数错误")
		return
	}

	qscer, err := models.FindQSCerByQscId(req.QscId)
	if err != nil {
		log.Errorf("Request Error: %s", err.Error())
		resp.Err(c, resp.WrongUsernameError, "找不到用户名")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(qscer.Password), []byte(req.Password))
	if err != nil {
		log.Errorf("Request Error: %s", err.Error())
		resp.Err(c, resp.WrongPasswordError, "密码错误")
		return
	}

	ss := sessions.Default(c)
	ss.Set(utils.SessionKeyUser, qscer)
	ss.Save()
	resp.Json(c, qscer)
}
