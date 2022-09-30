package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database //储存数据库
func Init() error {
	uri := "mongodb://127.0.0.1:27017"
	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetConnectTimeout(time.Second * 2)
	clientOptions.SetSocketTimeout(time.Second * 2)
	clientOptions.SetServerSelectionTimeout(time.Second * 2)

	Client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return err
	}

	err = Client.Ping(context.TODO(), nil)

	if err != nil {
		return err
	}

	DB = Client.Database("BBS")
	return nil
} //加载数据库
