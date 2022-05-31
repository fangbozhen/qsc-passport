package resp

import (
	. "passport-v4/global"

	"github.com/gin-gonic/gin"
)

// TODO: error code

type JsonResp struct {
	Err  string
	Data interface{}
}

func set_resp(c *gin.Context, resp JsonResp) {
	c.Set(CTX_RESPONSE, resp)
}

func JSON(c *gin.Context, obj interface{}) {
	set_resp(c, JsonResp{"", obj})
}

func ERR(c *gin.Context, str string) {
	set_resp(c, JsonResp{str, nil})
}
