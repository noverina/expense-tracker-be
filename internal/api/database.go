package api

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func InitDB() {
	timeout, err := strconv.Atoi(os.Getenv("CONNECT_TIMEOUT"))
	if err != nil {
		fmt.Errorf("unable to convert string to integer err=%w", err)
	}
	context, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	uri := os.Getenv("MONGODB_URI")
	client, err = mongo.Connect(context, options.Client().
		ApplyURI(uri))
	if err != nil {
		fmt.Errorf("unable to establish db connection err=%w", err)
	}
}

func Disconnect() {
	if client != nil {
		err := client.Disconnect(context.Background())
		if err != nil {
			LogError("unable to disconnect client", "err", err)
		}
	}
}

func GetDB() *mongo.Database {
	dbName := os.Getenv("MONGODB_DB")
	return client.Database(dbName)
}

func GetClient() *mongo.Client {
	return client
}
