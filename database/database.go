package database

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

type setting struct {
	DbUri              string `json:"dbUri"`
	Database           string `json:"database"`
	SessionExpiredTime int    `json:"sessionExpiredTime"`
}

var (
	Client                *mongo.Client
	s                     setting
	ctx                   = context.TODO()
	accountCountDocuments int64
)

type id struct {
	ID primitive.ObjectID `bson:"_id"`
}

func init() {
	file, err := os.ReadFile("setting.json")
	if err != nil {
		log.Panic(err)
		return
	}
	if err = json.Unmarshal(file, &s); err != nil {
		log.Panic(err)
	}

	clientOptions := options.Client().ApplyURI(s.DbUri).SetTimeout(10 * time.Second)
	// Connect to MongoDB
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	migrates(Account{}, Session{})
	countDocuments, err := Client.Database(s.Database).Collection("account").CountDocuments(ctx, bson.D{{"DeleteAt", false}})
	if err != nil {
		log.Println(err)
	}
	accountCountDocuments = countDocuments

	//for i := 20; i < 40; i++ {
	//	itoa := strconv.Itoa(i)
	//	a := NewAccount("user"+itoa, "user"+itoa)
	//	a.CreateAccount()
	//}
	fmt.Println("Connected to MongoDB!")
}

func migrates(interfaces ...interface{}) {
	type Collection struct {
		Name    string
		Options map[string]string
	}

	CollectionNames, err := Client.Database(s.Database).ListCollectionNames(ctx, map[any]any{}) // 取得 Collection 名稱清單
	if err != nil {
		log.Println("get ListCollectionNames Error", err)
	}

	inArray := func(s string, a []string) bool {
		for _, i := range a {
			if i == s {
				return true
			}
		}
		return false
	}

	for _, i := range interfaces {
		t := reflect.TypeOf(i)               // 先將接口進行解析
		ts := strings.Split(t.String(), ".") // 取得該接口的名稱
		tn := strings.ToLower(ts[len(ts)-1]) // 將接口名稱轉化為全小寫
		//fmt.Println("Collection Name: ", ts[len(ts)-1])

		if !inArray(tn, CollectionNames) { // 檢查 Collection 是否已存在
			err = Client.Database(s.Database).CreateCollection(ctx, tn) // 若該 Collection 不存在 就創建一個
			if err != nil {
				log.Println("CreateCollection Error:", err)
				log.Println("CollectionName: ", tn)
				return
			}
		}

		var indexModels []mongo.IndexModel
		for j := 0; j < t.NumField(); j++ {
			field := t.Field(j)
			//fmt.Println(field.Name, field.Tag.Get("op"))
			for _, tag := range strings.Split(field.Tag.Get("op"), ";") { // 解析自定義結構標籤
				switch strings.ToLower(tag) {
				case "unique":
					indexModels = append(indexModels, mongo.IndexModel{ // 新增一個 indexModel
						Keys:    bson.D{{field.Name, 1}},
						Options: options.Index().SetUnique(true),
					})
				}
			}

		}
		if len(indexModels) > 0 {
			_, err = Client.Database(s.Database).Collection(tn).Indexes().CreateMany(ctx, indexModels)
			if err != nil {
				log.Println(err)
			}
		} else {
			_, err = Client.Database(s.Database).Collection(tn).Indexes().DropAll(ctx)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func Disconnect() {
	err := Client.Disconnect(ctx)
	if err != nil {
		log.Println(err)
	}
}
