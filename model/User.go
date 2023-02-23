package model

import (
	"context"
	"errors"
	"fmt"
	"passport-v4/database"
	"passport-v4/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// enum LoginType
const (
	LoginZJU = "zju"
	LoginQSC = "qsc"
)

// enum position
const (
	PosIntern     = "实习成员"
	PosNormal     = "正式成员"
	PosConsultant = "顾问"
	PosManager    = "中管"
	PosMaster     = "高管"
)

// enum status
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
	Direction  string    `json:"direction" bson:"Direction"`
	Position   string    `json:"position" bson:"Position"`
	Status     string    `json:"status" bson:"Status"`
	Phone      string    `json:"phone" bson:"Phone"`
	Email      string    `json:"email" bson:"Email"`
	Note       string    `json:"note" bson:"Note"`
	Birthday   time.Time `json:"birthday,omitempty" bson:"Birthday"`
	JoinTime   time.Time `json:"jointime,omitempty" bson:"JoinTime"`
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

func checkUser(user *UserProfileQsc, DBUser *UserProfileQsc) {
	if user.Password == "" {
		user.Password = DBUser.Password
	}
	if user.Birthday.IsZero() {
		user.Birthday = DBUser.Birthday
	}
	if user.JoinTime.IsZero() {
		user.JoinTime = DBUser.JoinTime
	}
}

func UpdataOneByQscId(qscid string, user UserProfileQsc) error {
	var DBUser UserProfileQsc
	col := database.DB.Collection(utils.CollectionQscUsers)
	err := col.FindOne(ctx, bson.M{"QscId": qscid}).Decode(&DBUser)
	if err != nil {
		return err
	}
	checkUser(&user, &DBUser)
	res := col.FindOneAndReplace(ctx, bson.M{"QscId": qscid}, user)
	return res.Err()
}

func UpdateOne(qscid string, department string, position string) error {
	col := database.DB.Collection(utils.CollectionQscUsers)
	filter := bson.M{"QscId": qscid}
	var update bson.M
	if department == "" {
		update = bson.M{"$set": bson.M{"Position": position}}
	} else {
		update = bson.M{"$set": bson.M{"Department": department}}
	}
	_, err := col.UpdateMany(ctx, filter, update)
	return err
}

func InsertQSCer(user UserProfileQsc) error {
	col := database.DB.Collection(utils.CollectionQscUsers)
	_, err := col.InsertOne(ctx, user)
	return err
}

func DeleteByQscId(qscid string) error {
	col := database.DB.Collection(utils.CollectionQscUsers)
	res, err := col.DeleteOne(ctx, bson.M{"QscId": qscid})
	if res.DeletedCount == 0 {
		return errors.New(fmt.Sprintf("user %s not found", qscid))
	}
	return err
}

func FindInPages(selector interface{}, limit, page int64, sortCol string, isDescend bool) (users []UserProfileQsc, err error) {
	col := database.DB.Collection(utils.CollectionQscUsers)
	findOptions := options.Find()
	findOptions.SetSkip(page*limit - limit)
	findOptions.SetLimit(limit)
	if !isDescend {
		findOptions.SetSort(bson.D{{sortCol, -1}})
	} else {
		findOptions.SetSort(bson.D{{sortCol, 1}})
	}
	cur, err := col.Find(ctx, selector, findOptions)
	if err != nil {
		return nil, err
	}
	err = cur.All(ctx, &users)
	return
}
