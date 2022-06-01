package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"passport-v4/config"
	. "passport-v4/global"
	"passport-v4/model"
	"passport-v4/util/resp"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

/*** 当用户还没登录，跳转到登录界面 ***/

var zju_oauth *oauth2.Config

func ZJU_OauthInit() {
	zju_oauth = &oauth2.Config{
		ClientID:     config.ZjuOauth.ClientID,
		ClientSecret: config.ZjuOauth.ClientSecret,
		Scopes:       []string{}, // 看起来zju的oauth不需要这个
		Endpoint: oauth2.Endpoint{
			TokenURL:  "https://zjuam.zju.edu.cn/cas/oauth2.0/accessToken",
			AuthURL:   "https://zjuam.zju.edu.cn/cas/oauth2.0/authorize",
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: "https://qsc.zju.edu.cn/zju_oauth/code",
	}
}

func ZJU_LoginRequest(c *gin.Context) {

	var req struct {
		SuccessUrl string
		FailUrl    string
	}
	c.ShouldBind(&req)

	ss := sessions.Default(c)
	defer ss.Save()

	state := uuid.New().String()
	ss.Set(SS_KEY_STATE, state)
	ss.Set(SS_KEY_SUCCESS_URL, req.SuccessUrl)
	ss.Set(SS_KEY_FAILED_URL, req.FailUrl)

	url := zju_oauth.AuthCodeURL(state)
	c.Redirect(302, url)
}

func ZJU_OauthCodeReturn(c *gin.Context) {
	ctx := context.Background()
	ss := sessions.Default(c)
	defer ss.Save()
	code := c.Query("code")

	// 理论上需要判断回传state和session中的是否相等
	// 但是我看zjuam文档并没说返回这个值em
	// 所以拿不到这个值也就不报错了
	state := c.Query("state")
	if state != "" {
		session_state, ok := ss.Get(SS_KEY_STATE).(string)
		if !ok || session_state != state {
			resp.ERR(c, "state param incorrect")
		}
	}

	httpClient := &http.Client{Timeout: 2 * time.Second}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	/*** 后端通过code获取access_token ***/
	tok, err := zju_oauth.Exchange(ctx, code)
	if err != nil {
		redirect_login_failed(c, "cannot acquire access_token")
		return
	}

	ss.Set(SS_KEY_ACCESS_TOKEN, tok.AccessToken)

	zju_user, ok := get_zju_profile(tok.AccessToken)
	if !ok {
		redirect_login_failed(c, "cannot get zju profile")
		return
	}

	user := model.ZjuProfile2User(zju_user)

	ss.Set(SS_KEY_USER, user)

	redirect_login_success(c)
}

func redirect_login_failed(c *gin.Context, reason string) {
	ss := sessions.Default(c)
	url, ok := ss.Get(SS_KEY_FAILED_URL).(string)
	if !ok {
		logrus.Warn("login failed, but FAILED_URL not set")
		return
	}
	url = fmt.Sprintf("%s?reason=%s", url, reason)
	logrus.Infof("login failed: %s", reason)
	c.Redirect(302, url)
}

func redirect_login_success(c *gin.Context) {
	ss := sessions.Default(c)
	url, ok := ss.Get(SS_KEY_SUCCESS_URL).(string)
	if !ok {
		logrus.Warn("login success, but SUCCESS_URL not set")
		return
	}
	c.Redirect(302, url)
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
