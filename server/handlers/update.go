package handlers

import (
	"encoding/csv"
	"fmt"
	"passport-v4/model"
	"passport-v4/utils/resp"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// FIXME: 各种校验规则
func UpdateOne(c *gin.Context) {
	var req struct {
		Qscid string               `json:"qscid"`
		User  model.UserProfileQsc `json:"user"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Errorf("request error: %s", err.Error())
		resp.Err(c, resp.WrongRequestError, "参数错误")
		return
	}
	err = model.UpdataOneByQscId(req.Qscid, req.User)
	if err != nil {
		log.Errorf("update error: %s", err.Error())
		resp.Err(c, resp.DatabaseError, "数据库更新错误")
		return
	}
	log.Infof("sueccessfully update User: %s", req.Qscid)
	resp.Json(c, nil)
}

func UpdateMany(c *gin.Context) {
	var req struct {
		Ids        []string `json:"qscid"`
		Department string   `json:"department"`
		Position   string   `json:"position"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Errorf("request error: %s", err.Error())
		resp.Err(c, resp.WrongRequestError, "参数错误")
		return
	}
	for _, id := range req.Ids {
		err := model.UpdateOne(id, req.Department, req.Position)
		if err != nil {
			log.Errorf("err: %s", err.Error())
			resp.Err(c, resp.DatabaseError, "数据库批量更新失败")
			return
		}
	}
	log.Info("successfully update user: ", req.Ids)
	resp.Json(c, nil)
}

func Delete(c *gin.Context) {
	var req struct {
		Qscid string
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Errorf("request error: %s", err.Error())
		resp.Err(c, resp.WrongRequestError, "参数错误")
		return
	}
	err = model.DeleteByQscId(req.Qscid)
	if err != nil {
		log.Errorf("err: %s", err.Error())
		resp.Err(c, resp.DatabaseError, "数据库用户删除失败")
		return
	}
	log.Infof("successfully delete User: %s", req.Qscid)
	resp.Json(c, nil)
}

// FIXME: 各种校验规则
func Upload(c *gin.Context) {
	rfile, _ := c.FormFile("file")
	log.Infof("Get file: %s", rfile.Filename)
	file, _ := rfile.Open()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	record, _ := reader.ReadAll()

	mp := make(map[string]bool)
	var ids1, ids2 []string
	var flag bool
	flag = false
	for _, item := range record {

		if mp[item[1]] {
			ids1 = append(ids1, item[1])
			flag = true

		}
		if _, err := model.FindQSCerByQscId(item[1]); err == nil {
			ids2 = append(ids2, item[1])
			flag = true
		}
		mp[item[1]] = true
	}
	if flag {
		var reply string
		if len(ids1) != 0 {
			reply = fmt.Sprintf("用户表中QscId %s 重复  ", ids1)
		}
		if len(ids2) != 0 {
			reply = reply + fmt.Sprintf("数据库中已有 %s", ids2)
		}
		resp.Err(c, resp.DatabaseError, reply)
		return
	}

	for _, item := range record {
		pwd, _ := bcrypt.GenerateFromPassword([]byte(item[0]), bcrypt.DefaultCost)
		password := string(pwd)
		birthday, _ := time.Parse("2006/01/02", item[9])
		fmt.Println(birthday, item[9])
		user := model.UserProfileQsc{
			ZjuId:      item[0],
			QscId:      item[1],
			Password:   password,
			Name:       item[2],
			Gender:     item[3],
			Department: item[4],
			Position:   item[5],
			Status:     item[6],
			Phone:      item[7],
			Email:      item[8],
			Note:       "",
			Birthday:   birthday,
			JoinTime:   time.Now(),
		}
		err := model.InsertQSCer(user)
		if err != nil {
			log.Errorf("err: %s", err.Error())
			resp.Err(c, resp.DatabaseError, fmt.Sprintf("数据库用户%s插入失败", user.QscId))
			return
		}
	}

	resp.Json(c, nil)
}
