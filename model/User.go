package model

import (
	"context"
	"time"

	"passport-v4/database"
	. "passport-v4/global"

	"gopkg.in/mgo.v2/bson"
)

// enum LoginType
const (
	LT_ZJU = "zju"
	LT_QSC = "qsc"
)

// enum Position
const (
	POS_INTERN     = "实习成员"
	POS_NORMAL     = "正式成员"
	POS_CONSULTANT = "顾问"
	POS_MANAGER    = "中管"
	POS_MASTER     = "高管"
	POS_ADVANCED   = "高级成员"
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

// json是zjuam返回的，不能改
type UserProfileZju struct {
	Id         string  `json:"_id"`
	Attributes []ssmap `json:"attributes"`
}

type UserProfileQsc struct {
	Password   string    `json:"-" bson:"Password"` // hashed
	ZjuId      string    `json:"zjuid" bson:"ZjuId"`
	Name       string    `json:"name" bson:"Name"`
	QscId      string    `json:"qscid" bson:"QscId"`
	Gender     string    `json:"gender" bson:"Gender"`
	Position   string    `json:"position" bson:"Position"`     // 身份 @see enum Position
	Department string    `json:"department" bson:"Department"` // 部门
	Direction  string    `json:"direction" bson:"Direction"`   // 部门下分方向
	Status     string    `json:"status" bson:"Status"`         // 状态【保留】
	JoinTime   time.Time `json:"jointime" bson:"JoinTime"`     // 注意读出是GMT
	Privilege  smap      `json:"privilege" bson:"Privilege"`   // 权限组【保留】
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
		QscUser:   &pf,
	}
}

var ctx context.Context = context.TODO() //定义一个空的context,用于记录数据库操作的信息

func FindQSCerByQscId(qscid string) (UserProfileQsc, error) {
	col := database.DB.Collection(COL_QSC_USER)
	DBuser := UserProfileQsc{}
	err := col.FindOne(ctx, bson.M{"QscId": qscid}).Decode(&DBuser)
	return DBuser, err
}

func FindQSCerByZjuid(zjuid string) (UserProfileQsc, error) {
	col := database.DB.Collection(COL_QSC_USER)
	DBuser := UserProfileQsc{}
	err := col.FindOne(ctx, bson.M{"ZjuId": zjuid}).Decode(&DBuser)
	return DBuser, err
}

// 更改数据
func UpdateQSCer(user1 UserProfileQsc) error {
	col := database.DB.Collection(COL_QSC_USER)
	res := col.FindOneAndReplace(ctx, bson.M{"QscId": user1.QscId}, user1)
	return res.Err()
}

func InsertQSCer(user UserProfileQsc) error {
	col := database.DB.Collection(COL_QSC_USER)
	_, err := col.InsertOne(ctx, user)
	return err
}
