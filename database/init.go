package database

import "github.com/sirupsen/logrus"

func Init() (err error) {
	logrus.Info("[database]Init...")
	err = initMongo()
	if err != nil {
		return err
	}
	err = initRedis()
	if err != nil {
		return err
	}
	logrus.Info("[database] Init Success")
	return nil
}
