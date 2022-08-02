package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"passport-v4/config"
	. "passport-v4/global"
	"passport-v4/model"
	"passport-v4/util/resp"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

/*** 当用户还没登录，跳转到登录界面 ***/

func ZJU_LoginRequest(c *gin.Context) {

	var req struct {
		SuccessUrl string `form:"success"`
		FailUrl    string `form:"fail"`
	}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.ERR(c, resp.E_WRONG_REQUEST, "参数错误")
		return
	}

	ss := sessions.Default(c)

	ss.Set(SS_KEY_SUCCESS_URL, req.SuccessUrl)
	ss.Set(SS_KEY_FAILED_URL, req.FailUrl)
	ss.Save()

	url := fmt.Sprintf("https://zjuam.zju.edu.cn/cas/oauth2.0/authorize?client_id=%s&redirect_uri=%s&response_type=code",
		config.ZjuOauth.ClientID,
		url.QueryEscape(config.Server.UrlPrefix+"/zju/login_success"))

	// 302 会导致Cookie丢失
	c.HTML(200, "redirect.html", gin.H{
		"href": url,
	})
}

func ZJU_OauthCodeReturn(c *gin.Context) {
	ss := sessions.Default(c)
	code := c.Query("code")

	httpClient := &http.Client{Timeout: 2 * time.Second}
	url := fmt.Sprintf("https://zjuam.zju.edu.cn/cas/oauth2.0/accessToken?client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",
		config.ZjuOauth.ClientID,
		config.ZjuOauth.ClientSecret,
		code,
		url.QueryEscape(config.ZjuOauth.SsoUrl))
	r, err := httpClient.Get(url)

	if err != nil {
		fmt.Printf("err: %s", err.Error())
		redirect_login_failed(c, resp.E_INTERNAL_ERROR, "cannot acquire access_token")
		return
	}

	var tok struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		ErrCode     string `json:"errorcode"`
		ErrMsg      string `json:"errormsg"`
	}
	bs, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		redirect_login_failed(c, resp.E_INTERNAL_ERROR, "zjuam bad response")
		return
	}
	if bytes.HasPrefix(bs, []byte("error=")) {
		bs = bs[6:]
	}
	err = json.Unmarshal(bs, &tok)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		redirect_login_failed(c, resp.E_INTERNAL_ERROR, "zjuam bad response")
		return
	}
	if tok.ErrCode != "" {
		logrus.Errorf("zjuam failed %s %s", tok.ErrCode, tok.ErrMsg)

		redirect_login_failed(c, resp.E_INTERNAL_ERROR,
			"zjuam failed")
		return
	}
	ss.Set(SS_KEY_ACCESS_TOKEN, tok.AccessToken)
	ss.Save()
	zju_user, ok := get_zju_profile(tok.AccessToken)
	if !ok {
		redirect_login_failed(c, resp.E_INTERNAL_ERROR, "cannot get zju profile")
		return
	}

	user := model.ZjuProfile2User(zju_user)

	ss.Set(SS_KEY_USER, user)

	fmt.Printf("login success: %s %s", user.Name, user.ZjuId)

	redirect_login_success(c)
}

func redirect_login_failed(c *gin.Context, code int, reason string) {
	ss := sessions.Default(c)
	fmt.Printf("login failed: %d %s", code, reason)
	uri, ok := ss.Get(SS_KEY_FAILED_URL).(string)
	if !ok {
		logrus.Warn("login failed, but FAILED_URL not set")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	query := url.Values{}
	query.Set("reason", reason)
	query.Set("code", strconv.Itoa(code))
	uri = fmt.Sprintf("%s?%s", uri, query.Encode())
	logrus.Infof("login failed: [%d] %s", code, reason)

	c.Redirect(302, uri)
}

func redirect_login_success(c *gin.Context) {
	ss := sessions.Default(c)
	uri, ok := ss.Get(SS_KEY_SUCCESS_URL).(string)
	if !ok {
		logrus.Warn("login success, but SUCCESS_URL not set")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Redirect(302, uri)
}

func get_zju_profile(accass_token string) (user model.UserProfileZju, ok bool) {
	ok = false
	url := fmt.Sprintf("https://zjuam.zju.edu.cn/cas/oauth2.0/profile?access_token=%s", accass_token)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if err := json.Unmarshal(body, &user); err != nil {
		return
	}
	return user, true
}
