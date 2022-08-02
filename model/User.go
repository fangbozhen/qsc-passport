package model

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
	Id         string              `json:"id"`
	Attributes []map[string]string `json:"attributes"`
}

type UserProfileQsc struct {
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
	return User{LoginType: LT_QSC}
}
