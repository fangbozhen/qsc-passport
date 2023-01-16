package model

import "encoding/gob"

func Init() {
	gob.Register(UserProfileQsc{})
}
