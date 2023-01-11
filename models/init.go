package models

import "encoding/gob"

func Init() {
	gob.Register(UserProfileQsc{})
}
