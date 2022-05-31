package midware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// TODO: 还有POST BODY
func RemoveUndefinedQuery(c *gin.Context) {
	url := c.Request.URL
	url.RawQuery = strings.ReplaceAll(url.RawQuery, "undefined", "")
	url.RawQuery = strings.ReplaceAll(url.RawQuery, "null", "")
	c.Next()
}
