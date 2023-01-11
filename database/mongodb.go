package database

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var mgoCli *mongo.Client

func initDb() {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	clientOptions.SetConnectTimeout(time.Second * 2)
	clientOptions.SetConnectTimeout(time.Second * 2)
	clientOptions.SetServerSelectionTimeout(time.Second * 2)

	mgoCli, err = mongo.Connect(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Error while connecting to MongoDB: %s", err)
	}

	err = mgoCli.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Error in ping to mgoCli: %s", err)
	}
}

func MgoCli() *mongo.Client {
	if mgoCli == nil {
		initDb()
	}
	return mgoCli
}
