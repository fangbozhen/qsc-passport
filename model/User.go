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
	Name string `json:"NM"`
}

type UserProfileQsc struct {
}

func ZjuProfile2User(pf UserProfileZju) User {
	return User{LoginType: LT_ZJU}
}

func QscProfile2User(pf UserProfileQsc) User {
	return User{LoginType: LT_QSC}
}
