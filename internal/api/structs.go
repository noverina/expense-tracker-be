package api

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Dropdown struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Event struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Description string             `json:"description" bson:"description"`
	Type        string             `json:"type" bson:"type"`
	Category    string             `json:"category" bson:"category"`
	Date        time.Time          `json:"date" bson:"date"`
	Amount      string             `json:"amount" bson:"amount"`
}

type Category struct {
	Category string `bson:"category" json:"category"`
	Sum      string `bson:"sum" json:"sum"`
}

type Sum struct {
	Type       string     `bson:"type" json:"type"`
	Sum        string     `bson:"sum" json:"sum"`
	Categories []Category `bson:"categories" json:"categories"`
}

type Log struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	File     string             `json:"file" bson:"file"`
	Function string             `json:"function" bson:"function"`
	Error    string             `json:"error" bson:"error"`
	Date     time.Time          `json:"date" bson:"date"`
}

type Client struct {
	Identifier string `json:"identifier" bson:"_id"`
	SecretKey  string `json:"secret_key" bson:"secret_key"`
	Token      string `json:"token" bson:"token"`
	Role       string `json:"role" bson:"role"`
	Exp        int64  `json:"exp" bson:"exp"`
}

type Claims struct {
	jwt.RegisteredClaims
	Role string
}

type HttpResponse struct {
	IsError bool        `json:"is_error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Auth struct {
	Identifier string `json:"identifier"`
	SecretKey  string `json:"secret_key"`
}
