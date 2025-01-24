package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Category string             `json:"category" bson:"category"`
	Date     time.Time          `json:"date" bson:"date"`
	Amount   string 			`json:"amount" bson:"amount"`
}

var validFields = map[string]struct{}{
	"name":     {},
	"category": {},
	"date":     {},
	"amount":   {},
	"_id":      {},
}

func (e *Event) UnmarshalJSON(data []byte) error {
    type Alias Event
    inputJSON := &struct {
        ID string `json:"_id"`
		Date string `json:"date"`
		Amount string `json:"amount"`
        *Alias
    }{
        Alias: (*Alias)(e),
    }

    if err := json.Unmarshal(data, &inputJSON); err != nil {
        return errors.New("unable to unmarshal json into document")
    }

	if (!regexp.MustCompile(`^[0-9]*$`).MatchString(inputJSON.Amount)) {
		return errors.New("amount must only contain numbers")
	} else {
		e.Amount = inputJSON.Amount
	}

    if inputJSON.ID == "" {
        e.ID = primitive.NilObjectID
    } else {
        id, err := primitive.ObjectIDFromHex(inputJSON.ID)
        if err != nil {
            return errors.New("invalid objectId")
        }
        e.ID = id
    }

	if inputJSON.Date != "" {
		//TODO 2020-12-09T16:09:53+00:00
        parsedDate, err := time.Parse(time.RFC3339, inputJSON.Date)
        if err != nil {
            return errors.New("invalid date format")
        }
        e.Date = parsedDate
    }

    return nil
}

func UpsertEvent(c *gin.Context, event Event) (primitive.ObjectID, error) {
	if err := godotenv.Load(); err != nil {
		LogError("unable to load .env file", "err", err)
		return primitive.NilObjectID, err
	}
	collName := os.Getenv("MONGODB_COLL")
	coll := GetDB().Collection(collName);

	// no id, create new
	if (event.ID == primitive.NilObjectID) {
		event.ID = primitive.NewObjectID();
		result, err := coll.InsertOne(c, event) 
		if (err != nil) {
			LogError("unable to insert", "err", err, "id", event.ID)
			return primitive.NilObjectID, err
		}
		return result.InsertedID.(primitive.ObjectID), nil
	// id exist, update 
	} else {
		entity := bson.M{
        	"$set": bson.M{
            	"_id": event.ID,
				"name": event.Name,
				"category": event.Category,
				"date": event.Date,
				"amount": event.Amount,
        	},
    	}
		result, err := coll.UpdateByID(c, event.ID, entity)
		if (err != nil || result.ModifiedCount == 0) {
			LogError("unable to update", "err", err, "id", event.ID)
			return primitive.NilObjectID, err
		}
		return event.ID, nil
	}
}

func GetEventFilter(c *gin.Context, input map[string]interface{}) ([]Event, int, error) {
	collName := os.Getenv("MONGODB_COLL")
	coll := GetDB().Collection(collName);

	filter := bson.M{}
    for key, value := range input {
		if _, exists := validFields[key]; !exists {
			return nil, http.StatusBadRequest, errors.New("invalid field " + key)
		} else {
			filter[key] = value
		}
        
    }

	cursor, err := coll.Find(c, filter)
	if err != nil {
		LogError("unable to find document", "err", err)
		return nil, http.StatusInternalServerError, err
	}
	defer cursor.Close(c)

	var results []Event
	for cursor.Next(c) {
		var event Event
		if err := cursor.Decode(&event); err != nil {
			LogError("unable to decode cursor", "err", err)
			return nil, http.StatusInternalServerError, err
		}
		results = append(results, event)
	}
	if err := cursor.Err(); err != nil {
		LogError("cursor error", "err", err)
		return nil, http.StatusInternalServerError, err
	}
	if (len(results) == 0) {
		return nil, http.StatusBadRequest, errors.New("unable to find document with matching filter")
	}
	return results, http.StatusOK, nil
}