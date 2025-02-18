package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Event struct {
	ID       		primitive.ObjectID 	`json:"_id" bson:"_id"`
	Description     string      		`json:"description" bson:"description"`
	Type 			string             	`json:"type" bson:"type"`
	Category 		string             	`json:"category" bson:"category"`
	Date     		time.Time          	`json:"date" bson:"date"`
	Amount   		string 				`json:"amount" bson:"amount"`
}

var validFields = []string {}
var coll *mongo.Collection
var max int

func init () {
	validFields = append(validFields, "_id")
	validFields = append(validFields, "description")
	validFields = append(validFields, "type")
	validFields = append(validFields, "category")
	validFields = append(validFields, "date")
	validFields = append(validFields, "amount")
	collName := os.Getenv("MONGODB_COLL")
	coll = GetDB().Collection(collName);
	var err error
	max, err = strconv.Atoi(os.Getenv("MAX_EVENT_COUNT"))
	if (err != nil) {
		LogError("unable to convert string to integer")
	}
}

func (e *Event) UnmarshalJSON(data []byte) error {
    type Alias Event
    inputJSON := &struct {
        ID 			string `json:"_id"`
		Date 		string `json:"date"`
		Amount 		string `json:"amount"`
		Type 		string `json:"type"`
		Category 	string `json:"category"`
        *Alias
    }{
        Alias: (*Alias)(e),
    }

    if err := json.Unmarshal(data, &inputJSON); err != nil {
        return errors.New("unable to unmarshal json into document")
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

	if (!ValidType(inputJSON.Type)) {
		return errors.New("invalid type")
	}
	e.Type = inputJSON.Type

	if (e.Type == "income") {
		if (!ValidIncome(inputJSON.Category)) {
			return errors.New("invalid category")
		}
	} else {
		if (!ValidExpense(inputJSON.Category)) {
			return errors.New("invalid category")
		}
	}
	e.Category = inputJSON.Category

	if (!regexp.MustCompile(`^[0-9]*$`).MatchString(inputJSON.Amount)) {
		return errors.New("amount must only contain numbers")
	} 
	e.Amount = inputJSON.Amount
	

	//TODO 2020-12-09T16:09:53+00:00
    parsedDate, err := time.Parse(time.RFC3339, inputJSON.Date)
    if err != nil {
        return errors.New("invalid date format")
    }
    e.Date = parsedDate

	if (!regexp.MustCompile(`^[0-9]*$`).MatchString(inputJSON.Amount)) {
		return errors.New("amount must only contain numbers")
	}
	e.Amount = inputJSON.Amount
	

    return nil
}

func UpsertEvent(c *gin.Context, event Event) (primitive.ObjectID, int, error) {
	if err := godotenv.Load(); err != nil {
		LogError("unable to load .env file", "err", err)
		return primitive.NilObjectID, http.StatusInternalServerError, err
	}
	filter := make(map[string]interface{}) 
	var err error
	filter["date"] = event.Date
	exist, code, err := GetEventFilter(c, filter)
	if (len(exist) >= max) {
		return primitive.NilObjectID, http.StatusBadRequest, errors.New("event limit reached")
	} else if (code != 200 && err != nil) {
		return primitive.NilObjectID, code, err;
	}

	// no id, create new
	if (event.ID == primitive.NilObjectID) {
		event.ID = primitive.NewObjectID();
		result, err := coll.InsertOne(c, event) 
		if (err != nil) {
			LogError("unable to insert", "err", err, "id", event.ID)
			return primitive.NilObjectID, http.StatusInternalServerError, err
		}
		return result.InsertedID.(primitive.ObjectID), http.StatusOK, nil
	// id exist, update 
	} else {
		entity := bson.M{
        	"$set": bson.M{
            	"_id": event.ID,
				"type": event.Type,
				"category": event.Category,
				"description": event.Description,
				"date": event.Date,
				"amount": event.Amount,
        	},
    	}
		result, err := coll.UpdateByID(c, event.ID, entity)
		if (err != nil || result.ModifiedCount == 0) {
			LogError("unable to update", "err", err, "id", event.ID)
			return primitive.NilObjectID, http.StatusInternalServerError, err
		}
		return event.ID, http.StatusOK, nil
	}
}

func GetEventFilter(c *gin.Context, input map[string]interface{}) ([]Event, int, error) {
	filter := bson.M{}
    for key, value := range input {
		if exists := slices.Contains(validFields, key); !exists {
			return nil, http.StatusBadRequest, errors.New("invalid field " + key)
		} else {
			_, ok := value.(string)
			if (key == "date" && ok) {
				parsedDate, err := time.Parse(time.RFC3339, value.(string))
				if (err != nil) {
					return nil, http.StatusBadRequest, err
				}
				filter[key] = parsedDate
			} else {
				filter[key] = value
			}
			
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
	return results, http.StatusOK, nil
}

func GetEventByMonth(c *gin.Context, year string, month string, timezone string) ([]Event, int, error) {
	yearNum, err := strconv.Atoi(year)
	if (err != nil) {
		LogError("invalid year", "err", err)
		return nil, http.StatusBadRequest, err
	}
	monthNum, err := strconv.Atoi(month)
	if (err != nil) {
		LogError("invalid month", "err", err)
		return nil, http.StatusBadRequest, err
	}
	timezoneLoc, err := time.LoadLocation(timezone)
	if (err != nil) {
		LogError("invalid timezone", "err", err)
		return nil, http.StatusBadRequest, err
	}
	startDate := time.Date(yearNum, time.Month(monthNum), 1, 0, 0, 0, 0, timezoneLoc)
	lastDate := time.Date(yearNum, time.Month(monthNum) + 1, 1, 0, 0, 0, 0, timezoneLoc).AddDate(0, 0, -1)
	endDate := time.Date(yearNum, time.Month(monthNum), lastDate.Day(), 0, 0, 0, 0, timezoneLoc)

	filter := bson.M{"date": bson.M{
			"$gte": startDate, 
			"$lte": endDate,   
		}}

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
		event.Date = event.Date.In(timezoneLoc)
		results = append(results, event)
	}
	if err := cursor.Err(); err != nil {
		LogError("cursor error", "err", err)
		return nil, http.StatusInternalServerError, err
	}
	return results, http.StatusOK, nil
}