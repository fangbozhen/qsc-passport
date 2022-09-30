package database

import (
	"context"
	. "passport-v4/global"
	"passport-v4/model"

	"gopkg.in/mgo.v2/bson"
)

var ctx context.Context = context.TODO() //定义一个空的context,用于记录数据库操作的信息

func FindByQscId(user1 model.User) (model.UserProfileQsc, error) {
	col := DB.Collection(COL_USER)
	DBuser := model.UserProfileQsc{}
	err := col.FindOne(ctx, bson.M{"qscid": user1.Name}).Decode(&DBuser)
	return DBuser, err
}

func FindByName(user1 model.User) (model.UserProfileQsc, error) {
	col := DB.Collection(COL_USER)
	DBuser := model.UserProfileQsc{}
	err := col.FindOne(ctx, bson.M{"name": user1.Name}).Decode(&DBuser)
	return DBuser, err
}

// 更改数据
func UpdateUser(user1 model.UserProfileQsc) error {
	col := DB.Collection(COL_USER)
	res := col.FindOneAndReplace(ctx, bson.M{"id": user1.Id}, user1)
	return res.Err()
}

func Insert(user model.UserProfileQsc) error {
	col := DB.Collection("user")
	_, err := col.InsertOne(ctx, user)
	return err
}
