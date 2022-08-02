package resp

import (
	. "passport-v4/global"

	"github.com/gin-gonic/gin"
)

type JsonResp struct {
	Err  string      `json:"err"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func setResp(c *gin.Context, resp JsonResp) {
	c.Set(CTX_RESPONSE, resp)
}

func Json(c *gin.Context, obj interface{}) {
	setResp(c, JsonResp{"", 0, obj})
}

func Err(c *gin.Context, code int, str string) {
	setResp(c, JsonResp{str, code, nil})
}
