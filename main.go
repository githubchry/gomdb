package main

import (
	"github.com/githubchry/gomdb/drivers"
	"github.com/githubchry/gomdb/models"
	"github.com/githubchry/gomdb/testdata"
	"log"
	"time"
)

//https://www.cnblogs.com/Dr-wei/p/11742293.html
func main() {
	// log打印设置: Lshortfile文件名+行号  LstdFlags日期加时间
	log.SetFlags(log.Llongfile | log.LstdFlags | log.Lmicroseconds)

	var err error
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
	mgo := models.NewMdb("pedestrian")

	var p testdata.Pedestrian

	// 插入数据
	testdata.DeletePedestrianCollection(mgo)
	testdata.InsertRandomPedestrian(mgo, 1000, 2000);


	// 查询总数
	name, size := mgo.Count()
	log.Printf(" documents name: %+v documents size %d \n", name, size)
//*
	// 查询最后一个
	var timeStart time.Time


	timeStart = time.Now()
	singleResult := mgo.FindOne("eventid", 2222222)
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


//*/
/*
	//创建索引
	timeStart = time.Now()
	ret, err := mgo.CreateIndex("eventid")
	log.Printf("CreateIndex need %v\n", time.Since(timeStart))
	log.Print(ret, err)
*/
	// 断开连接
	drivers.MongoDBExit()
}
