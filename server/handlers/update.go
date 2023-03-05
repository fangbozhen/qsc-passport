package handlers

import (
	"fmt"
	"passport-v4/model"
	"passport-v4/utils/resp"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func checkUser(Department string, Direction string, Position string) (err error) {
	var (
		departments = [8]string{"新闻资讯中心", "产品研发中心", "人力资源部门", "摄影部", "推广策划中心", "视频部门", "设计与视觉中心", ""}
		directions  = [3]string{"技术研发方向", "产品运营方向", ""}
		positions   = [9]string{"实习成员", "正式成员", "中管", "顾问", "中级成员", "高管", "高级顾问", "退休老干部"}
	)
	flag := false
	for _, department := range departments {
		if Department == department {
			flag = true
		}
	}
	if !flag {
		return fmt.Errorf("部门字段错误")
	}
	if Department == departments[1] && Direction != directions[0] && Direction != directions[1] {
		return fmt.Errorf("方向字段错误")
	}
	flag = false
	for _, position := range positions {
		if Position == position {
			flag = true
		}
	}
	if !flag {
		return fmt.Errorf("职位字段错误")
	}
	return nil
}

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
	err = checkUser(req.User.Department, req.User.Direction, req.User.Position)
	if err != nil {
		log.Errorf("request error: %s", err.Error())
		resp.Err(c, resp.DatabaseError, err.Error())
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
		Direction  string   `json:"direction"`
		Position   string   `json:"position"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Errorf("request error: %s", err.Error())
		resp.Err(c, resp.WrongRequestError, "参数错误")
		return
	}
	err = checkUser(req.Department, req.Direction, req.Position)
	if err != nil {
		log.Errorf("request error: %s", err.Error())
		resp.Err(c, resp.DatabaseError, err.Error())
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

func checkRows(rows [][]string, c *gin.Context) bool {
	mp := make(map[string]bool)
	var ids1, ids2 []string
	var flag bool
	flag = false
	for _, row := range rows[1:] {
		if mp[row[9]] {
			ids1 = append(ids1, row[9])
			flag = true
		}
		if _, err := model.FindQSCerByQscId(row[9]); err == nil {
			ids2 = append(ids2, row[9])
			flag = true
		}
		mp[row[9]] = true
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
		return false
	}
	return true
}

func Upload(c *gin.Context) {
	rfile, err := c.FormFile("file")
	if err != nil {
		log.Errorf("Open file error: %s", err.Error())
		resp.Err(c, resp.WrongRequestError, "参数错误")
		return
	}
	file, err := rfile.Open()
	if err != nil {
		log.Errorf("Open file error: %s", err.Error())
		resp.Err(c, resp.WrongRequestError, "参数错误")
		return
	}
	defer file.Close()
	f, err := excelize.OpenReader(file)
	if err != nil {
		log.Errorf("Open file error: %s", err.Error())
		resp.Err(c, resp.WrongRequestError, "参数错误")
		return
	}
	defer f.Close()
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		log.Errorf("Open file error: %s", err.Error())
		resp.Err(c, resp.WrongRequestError, "参数错误")
		return
	}

	if !checkRows(rows, c) {
		return
	}

	for _, row := range rows[1:] {
		pwd, _ := bcrypt.GenerateFromPassword([]byte(row[1]), bcrypt.DefaultCost)
		password := string(pwd)
		birthday, _ := time.Parse("2006-01-02", row[3])
		user := model.UserProfileQsc{
			ZjuId:      row[1],
			QscId:      row[9],
			Password:   password,
			Name:       row[0],
			Gender:     row[2],
			Department: row[6],
			Direction:  row[7],
			Position:   row[8],
			Status:     model.StatusNormal,
			Phone:      row[5],
			Email:      row[4],
			Note:       "",
			Birthday:   birthday,
			JoinTime:   time.Now(),
			Privilege:  map[string]string{},
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
