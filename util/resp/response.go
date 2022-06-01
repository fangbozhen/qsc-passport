package resp

import (
	. "passport-v4/global"

	"github.com/gin-gonic/gin"
)

// TODO: error code

type JsonResp struct {
	Err  string
	Code int
	Data interface{}
}

func set_resp(c *gin.Context, resp JsonResp) {
	c.Set(CTX_RESPONSE, resp)
}

func JSON(c *gin.Context, obj interface{}) {
	set_resp(c, JsonResp{"", 0, obj})
}

func ERR(c *gin.Context, code int, str string) {
	set_resp(c, JsonResp{str, code, nil})
}
