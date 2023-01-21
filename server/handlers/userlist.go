package handlers

import (
	"QSCpassport/model"
	"QSCpassport/utils/resp"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetByPage(c *gin.Context) {
	type Sortby struct {
		Col     string `json:"col"`
		Boolean bool   `json:"isDescend"`
	}
	var req struct {
		PageNumber int64             `json:"pageNumber"`
		PageSize   int64             `json:"pageSize"`
		Sortby     Sortby            `json:"sortBy"`
		Filter     map[string]string `json:"filter"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Errorf("err: %s", err)
		resp.Err(c, resp.InternalError, "passport内部错误")
	}
	users, err := model.FindInPages(req.Filter, req.PageSize, req.PageNumber, req.Sortby.Boolean)
	if err != nil {
		log.Errorf("err: %s", err)
		resp.Err(c, resp.InternalError, "请检查输入参数")
	}
	resp.Json(c, gin.H{
		"users": users,
	})
}
