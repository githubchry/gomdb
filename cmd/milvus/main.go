package main

import (
	"encoding/json"
	"github.com/githubchry/gomdb/cmd"
	"github.com/githubchry/gomdb/drivers"
	"github.com/githubchry/gomdb/models"
	"github.com/milvus-io/milvus-sdk-go/milvus"
	"io/ioutil"
	"log"
	"os"
	"time"
)

//https://www.cnblogs.com/Dr-wei/p/11742293.html
func main() {
	// log打印设置: Lshortfile文件名+行号  LstdFlags日期加时间
	log.SetFlags(log.Llongfile | log.LstdFlags | log.Lmicroseconds)

	var err error

	//打开文件句柄操作
	fh, err := os.Open("E:\\ubuntu\\codes\\git\\PersonReIDServer\\release\\720.jpg")
	if err != nil {
		log.Fatal("error opening file")
	}
	defer fh.Close()
	img, _ := ioutil.ReadAll(fh)
	fh.Close()

	featureJson, err := cmd.GetFeatureData(img)
	log.Print(featureJson)
	/*
	{
	    "features":[
	        Array[2048],
	        Array[2048]
	    ],
	    "code":0,
	    "msg":"Success"
	}
	*/
	type featureJSON struct {
		Features [][]float32	`json:"features"`
		Code int	`json:"code"`
		Msg string	`json:"msg"`
	}
	var aa featureJSON
	timeStart := time.Now()
	err = json.Unmarshal([]byte(featureJson), &aa)
	timeElapsed := time.Since(timeStart)
	log.Println("Unmarshal", timeElapsed)

	log.Print(err)

	log.Print(aa.Msg)
	log.Print(len(aa.Features))
	log.Print(aa.Features[0][0])
	log.Print(aa.Features[0][6])



	//========================================
	milvusCfg := drivers.MilvusCfg{
		Addr:     "127.0.0.1",
		Port:     19530,
		Username: "chry",
		Password: "chry",
	}

	// 初始化连接到MongoDB
	err = drivers.MilvusDBInit(milvusCfg)
	if err != nil {
		log.Fatal(err)
	}

	//创建集合
	models.CreateCollection("person", 2048, 1024, int64(milvus.L2))

	//创建索引
	//models.CreateIndex("person", milvus.IVFSQ8)

	// 插入向量
	records := make([]milvus.Entity, len(aa.Features))
	for i := 0; i < len(aa.Features); i++ {
		records[i].FloatData = aa.Features[i]
	}
	id_array, err := models.InsertMany("person", "", records)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("id_array", id_array)
	}

	// 查询集合条数
	log.Print(models.Count("person"))

	// 打印集合信息
	models.GetCollectionInfo("person")

	//获取索引信息
	models.GetIndexInfo("person")

	// 以图搜图
	topkQueryResult, err := models.Search("person", records, 10)

	log.Println(len(topkQueryResult.QueryResultList))
	log.Printf("Search without index results: ")
	for i := 0; i < len(topkQueryResult.QueryResultList); i++ {
		for j := 0; j < len(topkQueryResult.QueryResultList[i].Ids); j++ {
			log.Print(topkQueryResult.QueryResultList[i].Ids[j], "=>", topkQueryResult.QueryResultList[i].Distances[j])
		}
	}

	// 断开连接
	drivers.MilvusDBExit()
}
