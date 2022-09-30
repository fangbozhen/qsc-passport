package database

import (
	"context"
	"passport-v4/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func Init() error {
	clientOptions := options.Client().ApplyURI(config.Mongo.Uri)
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

	DB = Client.Database(config.Mongo.Database)
	return nil
}
