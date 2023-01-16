package model

import (
	"QSCpassport/database"
	"QSCpassport/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// enum LoginType
const (
	LoginZJU = "zju"
	LoginQSC = "qsc"
)

//enum position
const (
	PosIntern     = "实习成员"
	PosNormal     = "正式成员"
	PosConsultant = "顾问"
	PosManager    = "中管"
	PosMaster     = "高管"
)

//enum status
const (
	StatusNormal = "在职"
	StatusRetire = "退休"
)

// 类型长了golang的自动对齐也太难看了

type smap map[string]interface{}
type ssmap map[string]string

type User struct {
	Name      string          `json:"Name"`
	ZjuId     string          `json:"ZjuId"`
	LoginType string          `json:"LoginType"` // @see enum LoginType
	QscUser   *UserProfileQsc `json:"QscUser,omitempty"`
}

// UserProfileZju json是zjuam返回的，不能改
type UserProfileZju struct {
	Id         string  `json:"_id"`
	Attributes []ssmap `json:"attributes"`
}

type UserProfileQsc struct {
	ZjuId      string    `json:"zjuid" bson:"ZjuId"`
	QscId      string    `json:"qscid" bson:"QscId"`
	Password   string    `json:"-" bson:"Password"`
	Name       string    `json:"name" bson:"Name"`
	Gender     string    `json:"gender" bson:"Gender"`
	Department string    `json:"department" bson:"Department"`
	Position   string    `json:"position" bson:"Position"`
	Status     string    `json:"status" bson:"Status"`
	Phone      string    `json:"phone" bson:"Phone"`
	Email      string    `json:"email" bson:"Email"`
	Birthday   time.Time `json:"birthday" bson:"Birthday"`
	JoinTime   time.Time `json:"jointime" bson:"JoinTime"`
}

func ZjuProfile2User(pf UserProfileZju) User {
	user := User{
		LoginType: LoginZJU,
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
		LoginType: LoginQSC,
		Name:      pf.Name,
		ZjuId:     pf.ZjuId,
		QscUser:   &pf,
	}
}

var ctx context.Context = context.TODO() //定义一个空的context,用于记录数据库操作的信息

func FindQSCerByQscId(qscid string) (UserProfileQsc, error) {
	col := database.DB.Collection(utils.CollectionQscUsers)
	DBuser := UserProfileQsc{}
	err := col.FindOne(ctx, bson.M{"QscId": qscid}).Decode(&DBuser)
	return DBuser, err
}

func FindQSCerByZjuid(zjuid string) (UserProfileQsc, error) {
	col := database.DB.Collection(utils.CollectionQscUsers)
	DBuser := UserProfileQsc{}
	err := col.FindOne(ctx, bson.M{"ZjuId": zjuid}).Decode(&DBuser)
	return DBuser, err
}

// 更改数据
func UpdateQSCer(user1 UserProfileQsc) error {
	col := database.DB.Collection(utils.CollectionQscUsers)
	res := col.FindOneAndReplace(ctx, bson.M{"QscId": user1.QscId}, user1)
	return res.Err()
}

func InsertQSCer(user UserProfileQsc) error {
	col := database.DB.Collection(utils.CollectionQscUsers)
	_, err := col.InsertOne(ctx, user)
	return err
}
