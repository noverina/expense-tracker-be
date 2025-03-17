package api

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Log struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	File     string             `json:"file" bson:"file"`
	Function string             `json:"function" bson:"function"`
	Error    string             `json:"error" bson:"error"`
	Date     time.Time          `json:"date" bson:"date"`
}

var logColl *mongo.Collection

func InitLog() {
	collName := os.Getenv("LOG_COLL")
	logColl = GetDB().Collection(collName)
}

func LogError(message string, keysAndValues ...interface{}) {
	pc, file, _, ok := runtime.Caller(1)
	if !ok {
		log.Error("unable to get function and file information")
		return
	}

	function := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	functionName := function[len(function)-1]
	fileSplit := strings.Split(file, "/")
	fileName := fileSplit[len(fileSplit)-1]

	entity := Log{
		ID:       primitive.NewObjectID(),
		File:     fileName,
		Function: functionName,
		Error:    message,
		Date:     time.Now(),
	}
	_, err := logColl.InsertOne(context.TODO(), entity)
	if err != nil {
		fmt.Errorf("unable to save log to db err=%w", err)
	}

	keysAndValues = append(keysAndValues, "function", functionName)
	keysAndValues = append(keysAndValues, "file", fileName)
	log.Error(message, keysAndValues...)
}

func LogWarn(message string, keysAndValues ...interface{}) {
	pc, file, _, ok := runtime.Caller(1)
	if !ok {
		log.Error("unable to get function and file information")
		return
	}

	function := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	functionName := function[len(function)-1]
	fileSplit := strings.Split(file, "/")
	fileName := fileSplit[len(fileSplit)-1]

	keysAndValues = append(keysAndValues, "function", functionName)
	keysAndValues = append(keysAndValues, "file", fileName)
	log.Warn(message, keysAndValues...)
}
