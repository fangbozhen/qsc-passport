package midware

import (
	"net/http"
	. "passport-v4/global"
	"passport-v4/util/resp"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Response(c *gin.Context) {
	c.Next()
	if c.Writer.Status() != 200 {
		return
	}
	obj, has := c.Get(CTX_RESPONSE)
	if !has {
		return
	}
	resp, ok := obj.(resp.JsonResp)
	if !ok {
		logrus.Errorf("[Response midware] CTX_RESPONSE Type Error")
		return
	}

	if resp.Err != "" {
		logrus.Infof("Err serve '%s': %s", c.Request.URL.Path, resp.Err)
	}

	c.JSON(http.StatusOK, resp)

}
