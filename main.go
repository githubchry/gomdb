package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/githubchry/gomdb/testdata"
	"log"

	"github.com/githubchry/gomdb/drivers"
	"github.com/githubchry/gomdb/models"
	"go.mongodb.org/mongo-driver/bson"
)

type Trainer struct {
	Name string
	Age  int
	City string
}
//https://www.cnblogs.com/Dr-wei/p/11742293.html
func main() {

	var p testdata.Pedestrian
	testdata.SetRandomPedestrian(&p)
	data, err := json.Marshal(p)
	if err != nil {
		log.Fatalln("json.Marshal error")
		return
	}
	fmt.Println("json:", string(data))

	//========================================
	cfg := drivers.MongoCfg{
		Addr:     "127.0.0.1",
		Port:     17017,		// 默认27017
		Username: "chry",
		Password: "chry",
		DBName:   "chrydb",
	}

	// 初始化连接到MongoDB
	err = drivers.MongoDBInit(cfg)
	if err != nil {
		log.Fatal(err)
	}

	mgo := models.NewMgo("trainers")
	// 单个插入
	ash := Trainer{"Ash", 10, "Pallet Town"}
	InsertOneResult := mgo.InsertOne(ash)
	fmt.Println("Inserted a single document: ", InsertOneResult)

	// 插入多个值
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Pewter City"}
	trainers := []interface{}{misty, brock}
	insertManyResult := mgo.InsertMany(trainers)
	fmt.Println("Inserted multiple documents: ", insertManyResult)

	// 更新
	update := bson.D{
		{"$inc", bson.D{
			{"age", 999},
		}},
	}
	updateResult := mgo.UpdateOne("name", "Ash", update)
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// 查询一个
	var result Trainer
	mgo.FindOne("name", "Ash").Decode(&result)
	fmt.Printf("Found a single document: %+v\n", result)

	// 查询总数
	name, size := mgo.Count()
	fmt.Printf(" documents name: %+v documents size %d \n", name, size)

	// 查询多个记录
	var results []*Trainer
	cur := mgo.FindAll(0, size, 1)
	defer cur.Close(context.TODO())
	if cur == nil {
		fmt.Println("FindAll err:", cur)
	}
	for cur.Next(context.TODO()) {
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// 遍历结果
	for k, v := range results {

		fmt.Printf("Found  documents  %d  %v \n", k, v)
	}

	// 删除文件
	deleteResult := mgo.DeleteMany("name", "Ash")
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult)

	// 断开连接
	drivers.MongoDBExit()
}
