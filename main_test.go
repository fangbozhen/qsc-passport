package main

import (
	"context"
	"fmt"
	"passport-v4/config"
	"passport-v4/database"
	"passport-v4/model"
	"passport-v4/utils"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func ResetPassword(qscid string) {
	col := database.DB.Collection(utils.CollectionQscUsers)
	cursor, _ := col.Find(
		context.Background(),
		// bson.M{"ZjuId": ""},
		bson.M{"QscId": qscid},
		options.Find().SetBatchSize(1000).SetProjection(bson.M{"QscId": 1, "ZjuId": 1}))
	var users []model.UserProfileQsc
	cursor.All(context.Background(), &users)

	for _, u := range users {

		fmt.Printf("reset pswd for %s\n", u.QscId)
		pswd, _ := bcrypt.GenerateFromPassword([]byte(u.ZjuId), bcrypt.MinCost)
		col.UpdateOne(
			context.Background(),
			bson.M{"QscId": u.QscId},
			bson.M{"$set": bson.M{"Password": string(pswd)}})
	}
}

func TestMain(t *testing.M) {
	config.Init()
	model.Init()
	database.InitDb()

	// ImportBbsUser2()
	// UpdateAdmin()
	// ResetPassword()
	// ImportNewFriends()
	ResetPassword("")
}
