package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func objectIDFromHex(hex string) primitive.ObjectID {
	objectID, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		log.Fatal(err)
	}
	return objectID
}
