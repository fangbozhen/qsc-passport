package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()

	/*** 当用户还没登录，跳转到登录界面 ***/

	conf := &oauth2.Config{
		ClientID:     "YOUR_CLIENT_ID",
		ClientSecret: "YOUR_CLIENT_SECRET",
		Scopes:       []string{}, // 看起来zju的oauth不需要这个
		Endpoint: oauth2.Endpoint{
			TokenURL:  "https://zjuam.zju.edu.cn/cas/oauth2.0/accessToken",
			AuthURL:   "https://zjuam.zju.edu.cn/cas/oauth2.0/authorize",
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: "https://qsc.zju.edu.cn/",
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	var user_state = ""
	url := conf.AuthCodeURL(user_state, oauth2.AccessTypeOffline)
	fmt.Printf("HTTP 302: %v", url)

	/*** 用户登录成功，向后端返回code ***/

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code = "[returned in oauth redirect]"

	// Use the custom HTTP client when requesting a token.
	httpClient := &http.Client{Timeout: 2 * time.Second}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	/*** 后端提供code获取access_token ***/

	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	/*** 后端拿着access_token去获取用户资料，完成session登录 ***/

	client := conf.Client(ctx, tok)
	_ = client
}
