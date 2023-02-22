package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"passport-v4/config"
	"passport-v4/model"
	"passport-v4/utils"
	"passport-v4/utils/resp"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/getsentry/sentry-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ZjuOauthRequest(c *gin.Context) {

	var req struct {
		SuccessUrl string `form:"success"`
		FailUrl    string `form:"fail"`
	}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		log.Errorf("err: %s", err.Error())
		resp.Err(c, resp.WrongRequestError, "参数错误")
		return
	}
	ss := sessions.Default(c)

	ss.Set(utils.SessionKeySuccessURL, req.SuccessUrl)
	ss.Set(utils.SessionKeyFailedURL, req.FailUrl)
	ss.Save()

	url := fmt.Sprintf("https://zjuam.zju.edu.cn/cas/oauth2.0/authorize?client_id=%s&redirect_uri=%s&response_type=code",
		config.ZjuOauth.ClientID,
		url.QueryEscape(config.Server.UrlPrefix+"/zju/login_success"))

	//  因为直接302会导致cookie保存失败，所以采用html内嵌js跳转
	c.String(200, "text/html;charset=utf-8", `<!DOCTYPE html>
	<html lang="en">
	<header>
		<meta charset="UTF-8">
		<title>重定向</title>
	</header>
	
	<body>
		<script>
			window.location.href = `+url+`;
		</script>
	</body>`)
}

func ZjuOauthCodeReturn(c *gin.Context) {
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
		log.Errorf("err: %s", err.Error())
		redirectLoginFailed(c, resp.InternalError, "cannot acquire access_token")
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
		log.Errorf("err: %s", err.Error())
		redirectLoginFailed(c, resp.InternalError, "zjuam bad response")
		return
	}
	if bytes.HasPrefix(bs, []byte("error=")) {
		bs = bs[6:]
	}
	err = json.Unmarshal(bs, &tok)
	if err != nil {
		log.Errorf("err: %s", err.Error())
		redirectLoginFailed(c, resp.InternalError, "zjuam bad response")
		return
	}
	if tok.ErrCode != "" {
		log.Errorf("zjuam failed %s %s", tok.ErrCode, tok.ErrMsg)

		redirectLoginFailed(c, resp.InternalError, "zjuam failed")
		return
	}
	ss.Set(utils.SessionKeyAccessToken, tok.AccessToken)
	zjuUser, ok := getZjuProfile(tok.AccessToken)
	if !ok {
		redirectLoginFailed(c, resp.InternalError, "cannot get zju profile")
		return
	}

	user := model.ZjuProfile2User(zjuUser)

	log.Infof("login success: %s %s", user.Name, user.ZjuId)

	// 如果是潮人，则继续登录为潮人
	if qscer, err := model.FindQSCerByZjuid(user.ZjuId); err == nil {
		log.Infof("Zjuer logined as QSCer: %s -> %s", user.Name, qscer.QscId)
		user = model.QscProfile2User(qscer)
	}

	ss.Set(utils.SessionKeyUser, user)
	ss.Save()

	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{ID: user.ZjuId, Name: user.Name})
		sentry.CaptureMessage("zju-login-success")
	})
	redirectLoginSuccess(c)
}

func redirectLoginFailed(c *gin.Context, code int, reason string) {
	ss := sessions.Default(c)
	log.Errorf("login failed: %d %s", code, reason)
	uri, ok := ss.Get(utils.SessionKeyFailedURL).(string)
	if !ok {
		log.Warn("login failed, but FAILED_URL not set")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	query := url.Values{}
	query.Set("reason", reason)
	query.Set("code", strconv.Itoa(code))
	uri = fmt.Sprintf("%s?%s", uri, query.Encode())

	log.Info("redirect to: %s", uri)
	c.Redirect(302, uri)
}

func redirectLoginSuccess(c *gin.Context) {
	ss := sessions.Default(c)
	uri, ok := ss.Get(utils.SessionKeySuccessURL).(string)
	if !ok {
		log.Warn("login success, but SUCCESS_URL not set")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	url, _ := url.Parse(uri)
	cookie, _ := c.Request.Cookie("SESSION_TOKEN")
	query := url.Query()
	query.Add("SESSION_TOKEN", cookie.Value)
	url.RawQuery = query.Encode()
	uri = url.String()

	log.Info("redirect to: %s", uri)
	c.Redirect(302, uri)
}

func getZjuProfile(accessToken string) (user model.UserProfileZju, ok bool) {
	ok = false
	url := fmt.Sprintf("https://zjuam.zju.edu.cn/cas/oauth2.0/profile?access_token=%s", accessToken)
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
