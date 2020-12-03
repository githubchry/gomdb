package main

import (
	"github.com/githubchry/gomdb/drivers"
	"github.com/githubchry/gomdb/models"
	"github.com/githubchry/gomdb/testdata"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

//https://www.cnblogs.com/Dr-wei/p/11742293.html
func main() {
	// log打印设置: Lshortfile文件名+行号  LstdFlags日期加时间
	log.SetFlags(log.Llongfile | log.LstdFlags | log.Lmicroseconds)

	var err error
	var timeStart time.Time

	//========================================
	mongoCfg := drivers.MongoCfg{
		Addr:     "127.0.0.1",
		Port:     17017,		// 默认27017
		Username: "chry",
		Password: "chry",
		DBName:   "chrydb",
	}

	// 初始化连接到MongoDB
	err = drivers.MongoDBInit(mongoCfg)
	if err != nil {
		log.Fatal(err)
	}
	mgo := models.NewMdb("pedestrian")

	var p testdata.Pedestrian

	// 查询总数
	size := mgo.Count()
	log.Printf("documents size %d \n", size)
	if size <= 1000000 {
		// 插入数据
		//testdata.DeletePedestrianCollection(mgo)
		//testdata.InsertRandomPedestrian(mgo, 1000, 2000);

		//创建索引
		//ret, err := mgo.CreateIndex("eventid")
		//log.Print(ret, err)
	}

	// =============== 删改查测试 ===============
	find100000 := func() {
		timeStart = time.Now()
		singleResult := mgo.FindOne("eventid", 100000)
		log.Printf("FindOne need %v\n", time.Since(timeStart))
		if singleResult.Err() != nil {
			log.Printf("%s\n", singleResult.Err())
		} else {
			err = singleResult.Decode(&p)
			if err != nil {
				log.Print(err)
			} else {
				log.Printf("%+v\n",p)
			}
		}
	}
	find100000()
	// 更新
	update := bson.D{
		{"$inc", bson.D{
			{"timestamp", 1},
		}},
	}
	updateResult := mgo.UpdateOne("eventid", 100000, update)
	log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	find100000()

	// 删除
	ret := mgo.DeleteOne("eventid", 100000)
	log.Printf("Delete %v document.\n", ret)
	find100000()

	// 插入
	mgo.InsertOne(p)
	find100000()

	// 断开连接
	drivers.MongoDBExit()
}
