package main

import (
	"encoding/json"
	"fmt"
	"os"
	"passport-v4/config"
	"passport-v4/database"
	"passport-v4/model"
	"testing"
	"time"

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

func TestMain(t *testing.M) {
	initLogger()
	Do(config.Init())
	Do(model.Init())
	Do(database.Init())

	ImportBbsUser()

}
