package model

import (
	"context"

	"passport-v4/database"
	. "passport-v4/global"

	"gopkg.in/mgo.v2/bson"
)

// enum LoginType
const (
	LT_ZJU = "zju"
	LT_QSC = "qsc"
)

type User struct {
	Name      string
	ZjuId     string
	LoginType string
}

type UserProfileZju struct {
	Id         string              `json:"_id"`
	Attributes []map[string]string `json:"attributes"`
}

type UserProfileQsc struct {
	Id         string                 `json:"_id"`
	PassWord   string                 `json:"password"`
	ZjuId      string                 `json:"zjuid"`
	Name       string                 `json:"name" bson:"name"`
	QscId      string                 `json:"qscid" bson:"qscid"`
	Gender     int                    `json:"gender"`
	Position   string                 `json:"position"`
	Department string                 `json:"department"`
	Status     int                    `json:"status"`
	Privilege  map[string]interface{} `json:"privilege"`
}

func ZjuProfile2User(pf UserProfileZju) User {
	user := User{
		LoginType: LT_ZJU,
		Name:      "",
		ZjuId:     "",
	}
	for _, item := range pf.Attributes {
		for k, v := range item {
			if k == "XM" {
				user.Name = v
			}
			if k == "CODE" {
				user.ZjuId = v
			}
		}
	}
	return user
}

func QscProfile2User(pf UserProfileQsc) User {
	return User{
		LoginType: LT_QSC,
		Name:      pf.Name,
		ZjuId:     pf.ZjuId,
	}
}

var ctx context.Context = context.TODO() //定义一个空的context,用于记录数据库操作的信息

func FindUserByQscId(user1 User) (UserProfileQsc, error) {
	col := database.DB.Collection(COL_QSC_USER)
	DBuser := UserProfileQsc{}
	err := col.FindOne(ctx, bson.M{"qscid": user1.Name}).Decode(&DBuser)
	return DBuser, err
}

func FindUserByName(user1 User) (UserProfileQsc, error) {
	col := database.DB.Collection(COL_QSC_USER)
	DBuser := UserProfileQsc{}
	err := col.FindOne(ctx, bson.M{"name": user1.Name}).Decode(&DBuser)
	return DBuser, err
}

// 更改数据
func UpdateUser(user1 UserProfileQsc) error {
	col := database.DB.Collection(COL_QSC_USER)
	res := col.FindOneAndReplace(ctx, bson.M{"id": user1.Id}, user1)
	return res.Err()
}

func InsertUser(user UserProfileQsc) error {
	col := database.DB.Collection(COL_QSC_USER)
	_, err := col.InsertOne(ctx, user)
	return err
}
