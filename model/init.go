package model

import "encoding/gob"

func Init() {
	gob.Register(User{})
	gob.Register(UserProfileQsc{})
	gob.Register(UserProfileZju{})
}
