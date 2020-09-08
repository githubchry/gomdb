package drivers

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDrivers struct {
	Client   *mongo.Client
	Database string
}

var MgoClient *mongo.Client
var MgoDbName string

// 初始化
func Init() {
	MgoClient = Connect()
	MgoDbName = "test"
}

// 连接
func Connect() *mongo.Client {

	// 设置客户端参数
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	//defer client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

// 关闭
func Close() {

	err := MgoClient.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
