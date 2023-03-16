package database

import (
	"crypto/sha256"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Account struct {
	ID         primitive.ObjectID `bson:"_id"`
	Username   string             `bson:"Username" json:"Username" op:"unique"`
	Password   string             `bson:"Password" json:"Password"`
	CreateTime primitive.DateTime `bson:"CreateTime"`
	UpdateTime primitive.DateTime `bson:"UpdateTime"`
	DeleteAt   bool               `bson:"DeleteAt"`
	Admin      bool               `bson:"Admin"`
}

func NewAccount(u, p string) *Account {
	h := sha256.New()
	h.Write([]byte(p + u))
	bs := h.Sum(nil)
	s := fmt.Sprintf("%x", bs)
	t := primitive.NewDateTimeFromTime(time.Now())
	return &Account{
		ID:         primitive.NewObjectID(),
		Username:   u,
		Password:   s,
		CreateTime: t,
		UpdateTime: t,
		DeleteAt:   false,
	}
}

// GetPasswordHash 適用於密碼未加密的狀態 回傳加密後的密碼
func (a *Account) GetPasswordHash() string {
	h := sha256.New()
	h.Write([]byte(a.Password + a.Username))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
func (a *Account) CreateAccount() error {
	if _, err := Client.Database(s.Database).Collection("account").InsertOne(ctx, a); err != nil {
		return err
	}
	return nil
}

func (a *Account) GetID() string {
	SingleResult := Client.Database(s.Database).Collection("account").FindOne(ctx, a)
	var id id
	if err := SingleResult.Decode(&id); err != nil {
		log.Println(err)
	} else {
		return id.ID.Hex()
	}
	return ""
}

// Update 更新用戶資訊
func (a *Account) Update() error {
	coll := Client.Database(s.Database).Collection("account")
	update := bson.D{
		{"$set", bson.D{
			{"Username", a.Username},
			{"Password", a.Password},
			{"Admin", a.Admin},
			{"UpdateTime", primitive.NewDateTimeFromTime(time.Now())},
		}},
	}
	_, err := coll.UpdateByID(ctx, a.ID, update)
	//err = errors.New("測試")
	return err
}

func (a *Account) Delete() error {
	coll := Client.Database(s.Database).Collection("account")
	update := bson.D{
		{"$set", bson.D{
			{"DeleteAt", true},
		}},
	}
	_, err := coll.UpdateByID(ctx, a.ID, update)
	return err
}

// SearchAccount 給Username 查詢用戶
func SearchAccount(username string) (*Account, error) {
	var a Account
	SingleResult := Client.Database(s.Database).Collection("account").FindOne(ctx, bson.D{{"Username", username}})
	if err := SingleResult.Decode(&a); err != nil {
		return nil, err
	}
	return &a, nil
}

// GetAccountFromSession 利用 Session 找出用戶
func GetAccountFromSession(sess string) (a *Account, err error) {
	one := Client.Database(s.Database).Collection("session").FindOne(ctx, bson.D{
		{"Session", sess},
		{"ExpiredTime", bson.D{
			{"$gte", time.Now()},
		}},
	})
	var se Session
	if err = one.Decode(&se); err != nil {
		log.Println("Session 格式化錯誤", err)
		return
	}
	coll := Client.Database(s.Database).Collection("account")
	c := coll.FindOne(ctx, bson.D{{"_id", objectIDFromHex(se.AccountID)}})
	if err = c.Decode(&a); err != nil {
		log.Println("Account 格式化錯誤", err)
		return
	}
	return
}

// GetAccountList 取得用戶清單
func GetAccountList(limit, page int64) (accounts []*Account, err error) {
	opts := options.Find().SetLimit(limit).SetSkip(limit * page)
	coll := Client.Database(s.Database).Collection("account")
	find, err := coll.Find(ctx, bson.D{{"DeleteAt", false}}, opts)
	//find, err := coll.Find(ctx, bson.D{{"DeleteAt", false}})
	if err != nil {
		return nil, err
	}
	err = find.All(ctx, &accounts)
	return
}

// GetAccountCountDocuments 取得用戶總量
func GetAccountCountDocuments() int64 {
	return accountCountDocuments
}

// ReGetAccountCountDocuments 取得並更新用戶總量
func ReGetAccountCountDocuments() int64 {
	countDocuments, err := Client.Database(s.Database).Collection("account").CountDocuments(ctx, bson.D{{"DeleteAt", false}})
	if err != nil {
		log.Println(err)
	}
	accountCountDocuments = countDocuments
	return accountCountDocuments
}

func GetAccountFromID(id primitive.ObjectID) (a *Account, err error) {
	one := Client.Database(s.Database).Collection("account").FindOne(ctx, bson.D{{"_id", id}})
	err = one.Decode(&a)
	return
}
