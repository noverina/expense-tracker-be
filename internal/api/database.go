package api

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Init() {
	if err := godotenv.Load(); err != nil {
		fmt.Errorf("unable to load .env file err=%w", err)
	}

	timeout, err := strconv.Atoi(os.Getenv("CONNECT_TIMEOUT"))
	if (err != nil) {
		fmt.Errorf("unable to load .env file err=%w", err)
	}
	context, cancel := context.WithTimeout(context.Background(), time.Duration(timeout) * time.Second)
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
			LogError("database.go", "Disconnect", "unable to disconnect client", "err", err)
		} 
	}
}

func GetDB() *mongo.Database {
	dbName := os.Getenv("MONGODB_DB")
	return client.Database(dbName);
}

func GetClient() *mongo.Client {
	return client
}

func ConvertStrToObjId(id string) primitive.ObjectID  {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		LogError("database.go", "ConvertStrToObjId", err.Error())
	}
	return objId
}