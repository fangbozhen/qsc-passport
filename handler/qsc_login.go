package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	. "passport-v4/global"
	"passport-v4/model"
	"passport-v4/util/resp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const qsc_login_url = "qsc.zju.edu.cn/lllooogggiiinnn"

func QSC_Login(c *gin.Context) {

	var req struct {
		Username string
		Password string
	}
	var qsc_user model.UserProfileQsc
	var user model.User

	ss := sessions.Default(c)
	defer ss.Save()

	query := url.Values{}
	query.Set("username", req.Username)
	query.Set("password", req.Password)

	rs, err := http.Get(fmt.Sprintf("%s?%s", qsc_login_url, query.Encode()))
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return
	}
	if err := json.Unmarshal(body, &qsc_user); err != nil {
		return
	}
	// TODO 判断失败情况

	user = model.QscProfile2User(qsc_user)
	ss.Set(SS_KEY_USER, user)

	resp.JSON(c, user)

}
