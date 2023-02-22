package middleware

import (
	"net/http"
	"passport-v4/utils"
	"passport-v4/utils/resp"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Response(c *gin.Context) {
	c.Next()
	if c.Writer.Status() != 200 {
		return
	}
	obj, has := c.Get(utils.CtxResponse)
	if !has {
		return
	}
	resp, ok := obj.(resp.JsonResp)
	if !ok {
		log.Errorf("[Response midware] CTX_RESPONSE Type Error")
		return
	}

	if resp.Err != "" {
		log.Infof("Err serve '%s': %s", c.Request.URL.Path, resp.Err)
	}

	c.JSON(http.StatusOK, resp)

}
