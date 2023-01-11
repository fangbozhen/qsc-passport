package models

import (
	"QSCpassport/database"
	"QSCpassport/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

//enum position
const (
	POS_INTERN     = "实习成员"
	POS_NORMAL     = "正式成员"
	POS_CONSULTANT = "顾问"
	POS_MANAGER    = "中管"
	POS_MASTER     = "高管"
)

//enum status
const (
	STATUS_NORMAL = "在职"
	STATUS_RETIRE = "退休"
)

type (
	// 数据表表结构体类
	UserProfileQsc struct {
		ZjuId      string    `json:"zjuid" bson:"ZjuId"`
		QscId      string    `json:"qscid" bson:"QscId"`
		Password   string    `json:"-" bson:"Password"`
		Name       string    `json:"name" bson:"Name"`
		Gender     string    `json:"gender" bson:"Gender"`
		Department string    `json:"department" bson:"Department"`
		Direction  string    `json:"direction,omitempty" bson:"Direction"`
		Position   string    `json:"position" bson:"Position"`
		Status     string    `json:"status" bson:"Status"`
		Phone      string    `json:"phone" bson:"Phone"`
		Email      string    `json:"email" bson:"Email"`
		Birthday   time.Time `json:"birthday" bson:"Birthday"`
		JoinTime   time.Time `json:"jointime" bson:"JoinTime"`
	}

	// 请求结构体类
	UserReq struct {
		QscId    string `json:"qscid"`
		Password string `json:"password"`
	}
)

func getCollection(db string, collection string) *mongo.Collection {
	client := database.MgoCli()
	return client.Database(db).Collection(collection)
}

func FindQSCerByQscId(qscid string) (UserProfileQsc, error) {
	collection := getCollection("my_db", utils.CollectionQscUsers)
	user := UserProfileQsc{}
	err := collection.FindOne(context.TODO(), bson.M{"QscId": qscid}).Decode(&user)
	return user, err
}
