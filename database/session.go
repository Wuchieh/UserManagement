package database

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type Session struct {
	Session     string             `bson:"Session" op:"unique"`
	AccountID   string             `bson:"Account"`
	ExpiredTime primitive.DateTime `bson:"ExpiredTime"`
}

func NewSession(s, a string) *Session {
	return &Session{
		Session:   s,
		AccountID: a,
	}
}

func (se *Session) CreateSession(MaxAge ...int) error {
	//se.ExpiredTime = primitive.NewDateTimeFromTime(time.Now().Add(time.Duration(int64(s.SessionExpiredTime)) * time.Second))
	if len(MaxAge) > 0 && MaxAge[0] > 0 {
		se.ExpiredTime = primitive.NewDateTimeFromTime(time.Now().Add(time.Duration(MaxAge[0]) * time.Second))
	} else {
		se.ExpiredTime = primitive.NewDateTimeFromTime(time.Now().Add(time.Duration(s.SessionExpiredTime) * time.Second))
	}
	_, err := Client.Database(s.Database).Collection("session").InsertOne(ctx, se)
	return err
}

func (se *Session) Drop() {
	result, err := Client.Database(s.Database).Collection("session").DeleteOne(ctx, se)
	if err != nil {
		return
	}
	fmt.Printf("Number of documents deleted: %d\n", result.DeletedCount)
}

func SearchSession(sess string) (*Session, error) {
	var se Session
	SingleResult := Client.Database(s.Database).Collection("session").FindOne(ctx, bson.D{{"Session", sess}})
	if err := SingleResult.Decode(&se); err != nil {
		return nil, err
	}
	return &se, nil
}

func ClearSessions(expiredSessionList ...string) error {
	filterA := bson.A{}
	for _, s := range expiredSessionList {
		filterA = append(filterA, bson.D{{"Session", s}})
	}
	filter := bson.D{{"$or", filterA}}
	many, err := Client.Database(s.Database).Collection("session").DeleteMany(ctx, filter)
	fmt.Printf("刪除了 %d 個Session\n", many.DeletedCount)
	return err
}

func ClearExpiredSession() {
	coll := Client.Database(s.Database).Collection("session")
	finds, err := coll.DeleteMany(ctx, bson.D{{"ExpiredTime", bson.D{{"$lt", time.Now()}}}})
	if err != nil {
		log.Fatal("Error: ", err)
		return
	}
	if finds.DeletedCount != 0 {
		log.Println("已刪除了 ", finds.DeletedCount, "個過期的 Session")
	}
}
