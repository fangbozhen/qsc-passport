package model

import "encoding/gob"

// 让redis认识我们的结构体
func Init() error {
	gob.Register(User{})
	return nil
}
