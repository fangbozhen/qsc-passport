package handlers

import (
	"passport-v4/model"
	"passport-v4/utils/resp"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// TODO: isDescend 是不是反了
func GetByPage(c *gin.Context) {
	type Sortby struct {
		Col     string `json:"col"`
		Boolean bool   `json:"isDescend"`
	}
	var req struct {
		PageNumber int64             `json:"pageNumber"`
		PageSize   int64             `json:"pageSize"`
		Filter     map[string]string `json:"filter"`
		Sortby     Sortby            `json:"sortBy"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Errorf("err: %s", err)
		resp.Err(c, resp.WrongRequestError, "参数错误")
		return
	}
	users, err := model.FindInPages(req.Filter, req.PageSize, req.PageNumber, req.Sortby.Col, req.Sortby.Boolean)
	if err != nil {
		log.Errorf("err: %s", err)
		resp.Err(c, resp.DatabaseError, "数据库错误")
		return
	}
	resp.Json(c, gin.H{
		"users": users,
	})
}
