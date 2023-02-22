package database

import (
	"context"
	"passport-v4/config"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func InitDb() {
	log.Info("[Database] Init...")
	clientOptions := options.Client().ApplyURI(config.Mongo.Uri)
	clientOptions.SetConnectTimeout(time.Second * 2)
	clientOptions.SetConnectTimeout(time.Second * 2)
	clientOptions.SetServerSelectionTimeout(time.Second * 2)

	mgoCli, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Error while connecting to MongoDB: %s", err)
	}

	err = mgoCli.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Error in ping to mgoCli: %s", err)
	}

	DB = mgoCli.Database(config.Mongo.Database)
	log.Info("[Database] Init success")
}
