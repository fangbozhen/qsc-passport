package handlers

import (
	"context"
	"passport-v4/database"
	"passport-v4/utils"
	"passport-v4/utils/resp"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetPivilege(c *gin.Context) {
	col := database.DB.Collection(utils.CollectionList)
	var res struct {
		Name      string   `bson:"Name"`
		Privilege []string `bson:"Privilege"`
	}
	err := col.FindOne(context.TODO(), bson.M{"Name": "privilege"}).Decode(&res)
	if err != nil {
		log.Errorf("database error: %s", err.Error())
		resp.Err(c, resp.DatabaseError, "数据库错误")
		return
	}
	resp.Json(c, gin.H{
		"privilege": res.Privilege,
	})
}
