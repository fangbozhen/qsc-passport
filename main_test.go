package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"passport-v4/config"
	"passport-v4/database"
	"passport-v4/global"
	"passport-v4/model"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func ImportBbsUser() {

	var users []struct {
		Qscid      string `json:"discuz_username"`
		Name       string `json:"name"`
		Stuid      string `json:"stuid"`
		JoinTime   string `json:"time"`
		Department string `json:"department"`
	}

	f, _ := os.Open("qscbbs_passport.json")
	json.NewDecoder(f).Decode(&users)

	for _, u := range users {
		t, _ := time.Parse("2006-01-02 15:04:05", u.JoinTime)
		if u.Name == "" {
			u.Name = u.Qscid
		}
		if u.Department == "" {
			u.Department = "老老人"
		}

		qscer := model.UserProfileQsc{
			QscId:      u.Qscid,
			Password:   u.Qscid,
			ZjuId:      u.Stuid,
			Name:       u.Name,
			Department: u.Department,
			JoinTime:   t,
			Position:   model.POS_NORMAL,
			Gender:     "",
			Status:     "",
			Privilege:  map[string]interface{}{},
		}
		hashed, _ := bcrypt.GenerateFromPassword([]byte(qscer.Password), bcrypt.MinCost)
		qscer.Password = string(hashed)
		model.InsertQSCer(qscer)
		fmt.Printf("Welcome %s %s\n", qscer.Name, qscer.QscId)
	}
}

func ImportBbsUser2() {

	var users []struct {
		Qscid      string `json:"username"`
		Name       string
		Stuid      string
		JoinTime   int64 `json:"regdate"`
		Department string
	}

	f, _ := os.Open("qsc_bbs_common_member.json")
	json.NewDecoder(f).Decode(&users)

	for _, u := range users {
		t := time.Unix(u.JoinTime, 0)
		if cnt, _ := database.DB.Collection(global.COL_QSC_USER).CountDocuments(context.Background(), bson.M{"QscId": u.Qscid}); cnt != 0 {
			fmt.Printf("dup qscid: %s    passed\n", u.Qscid)
			continue
		}
		if u.Name == "" {
			u.Name = u.Qscid
		}
		if u.Department == "" {
			u.Department = "老老人"
		}
		qscer := model.UserProfileQsc{
			QscId:      u.Qscid,
			Password:   u.Qscid,
			ZjuId:      u.Stuid,
			Name:       u.Name,
			Department: u.Department,
			JoinTime:   t,
			Position:   model.POS_NORMAL,
			Gender:     "",
			Status:     "",
			Privilege:  map[string]interface{}{},
		}
		hashed, _ := bcrypt.GenerateFromPassword([]byte(qscer.Password), bcrypt.MinCost)
		qscer.Password = string(hashed)
		model.InsertQSCer(qscer)
		fmt.Printf("Welcome %s %s\n", qscer.Name, qscer.QscId)
	}
}

func UpdateAdmin() {
	file, _ := os.Open("admin.csv")
	reader := csv.NewReader(file)
	rows, _ := reader.ReadAll()
	for _, row := range rows[1:] {
		var (
			qscid = row[1]
			name  = row[2]
		)
		res, err := database.DB.Collection(global.COL_QSC_USER).UpdateOne(
			context.Background(),
			bson.M{"$or": bson.A{
				bson.M{"Name": name},
				bson.M{"QscId": qscid},
			}},
			bson.M{
				"$set": bson.M{"Position": model.POS_MANAGER},
			})
		if err != nil || res.ModifiedCount != 1 {
			fmt.Printf("bad qscid: %s\n", qscid)
		}
		// fmt.Printf("update: %s", qscid)
	}
}

func SlContains(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func ImportNewFriends() {
	col := database.DB.Collection(global.COL_QSC_USER)
	file, _ := os.Open("./newbie.csv")
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Errorf("%v", err)
	}
	for _, row := range rows {
		if row[0] != "实习成员" {
			continue
		}
		var (
			name   = row[1]
			gender = row[2]
			qscid  = row[3]
			zjuid  = row[5]
			dept   = row[10]
		)

		if SlContains([]string{"浪味仙", "苏苏", "夏至", "瓜瓜", "洛洛"}, qscid) {
			qscid = qscid + "2"
		} else {
			continue
		}

		if cnt, _ := col.CountDocuments(context.Background(), bson.M{"QscId": qscid}); cnt != 0 {
			fmt.Printf("%s is same qscid!\n", qscid)
			continue
		}

		pswd, _ := bcrypt.GenerateFromPassword([]byte(zjuid), bcrypt.MinCost)

		direction := ""
		if dept == "技术研发中心" {
			direction = "技术研发方向"
			dept = "产品研发中心"
		}
		if dept == "产品运营部门" {
			direction = "产品运营方向"
			dept = "产品研发中心"
		}

		col.InsertOne(
			context.Background(),
			model.UserProfileQsc{
				ZjuId:      zjuid,
				QscId:      qscid,
				Name:       name,
				Password:   string(pswd),
				Gender:     gender,
				Position:   model.POS_INTERN,
				Department: dept,
				Direction:  direction,
				JoinTime:   time.Now(),
			})
		fmt.Printf("Welcome %s %s\n", name, qscid)
	}
}

func ResetPassword(qscid string) {
	cursor, _ := database.DB.Collection(global.COL_QSC_USER).Find(
		context.Background(),
		// bson.M{"ZjuId": bson.M{"$ne": ""}},
		bson.M{"QscId": qscid},
		options.Find().SetBatchSize(1000).SetProjection(bson.M{"QscId": 1, "ZjuId": 1}))
	var users []model.UserProfileQsc
	cursor.All(context.Background(), &users)

	for _, u := range users {

		fmt.Printf("reset pswd for %s\n", u.QscId)
		pswd, _ := bcrypt.GenerateFromPassword([]byte(u.ZjuId), bcrypt.MinCost)
		database.DB.Collection(global.COL_QSC_USER).UpdateOne(
			context.Background(),
			bson.M{"QscId": u.QscId},
			bson.M{"$set": bson.M{"Password": string(pswd)}})
	}
}

func TestMain(t *testing.M) {
	initLogger()
	Do(config.Init())
	Do(model.Init())
	Do(database.Init())

	// ImportBbsUser2()
	// UpdateAdmin()
	// ResetPassword()
	// ImportNewFriends()
	// ResetPassword("创高")
}
