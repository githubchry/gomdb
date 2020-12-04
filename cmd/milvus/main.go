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
	"strconv"
	"time"
)

//https://www.cnblogs.com/Dr-wei/p/11742293.html
func main() {
	// log打印设置: Lshortfile文件名+行号  LstdFlags日期加时间
	log.SetFlags(log.Llongfile | log.LstdFlags | log.Lmicroseconds)
	//log.Println("go","python","php","javascript") // go python php javascript
	//log.Print("go","python","php","javascript") // gopythonphpjavascript
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



	//================================================================
	milvusCfg := drivers.MilvusCfg{
		Addr:     "127.0.0.1",
		//Addr:     "10.2.3.5",
		Port:     19530,
		Username: "chry",
		Password: "chry",
	}

	// 初始化连接到Milvus
	err = drivers.MilvusDBInit(milvusCfg)
	if err != nil {
		log.Fatal(err)
	}

	//================================================================

	//获取配置信息
	cfg, err := models.GetConfig()
	log.Print(cfg)

	//获取集合信息
	collections, err := models.ListCollections()
	for i := 0; i < len(collections); i++ {
		// 打印集合信息
		dimension, indexFileSize, _ := models.GetCollectionInfo(collections[i])
		partitionNames, _ := models.ListPartitions(collections[i])
		indexParam, _ := models.GetIndexInfo(collections[i])
		log.Printf("%d维集合[%s], 索引大小[%dM], 索引类型[%v], 索引参数[%s], 当前共[%d]条记录, 有[%d]个分区:%v", dimension, collections[i], indexFileSize, indexParam.IndexType, indexParam.ExtraParams, models.Count(collections[i]), len(partitionNames), partitionNames)
	}

	//================================================================

	//删除集合 此后操作必须等1秒以上
	models.DropCollection("person" )
	time.Sleep(time.Second*1)

	//创建集合
	models.CreateCollection("person", 2048, 1024, int64(milvus.IP))

	// 删除索引
	models.DropIndex("person")

	//创建索引
	indexParam := milvus.IndexParam{
		CollectionName: "person",
		IndexType:      milvus.IVFFLAT,
		ExtraParams:    "{\"nlist\" : 16384, \"nprob\" : 32}",
	}
	models.CreateIndex(indexParam)

	// 创建分区, 不能超过4096个
	for i := 0; i < 5; i++ {
		partName := "part" + strconv.Itoa(i)
		models.CreatePartition("person", partName)
		log.Println("Create collection : " +partName)
	}

	//下面会执行失败, 因为Default partition cannot be dropped.
	models.DropPartition("person", "_default")
	//================================================================

	// 插入向量
	records := make([]milvus.Entity, len(aa.Features))
	for i := 0; i < len(aa.Features); i++ {
		records[i].FloatData = aa.Features[i]
	}
	id_array, err := models.Insert("person", "", records)
	if err != nil {
		log.Print(err)
	} else {
		log.Print("id_array", id_array)
	}

	//预加载集合  搜索前加载可提高速度
	models.LoadCollection("person")

	// 以图搜图
	topkQueryResult, err := models.Search("person", nil, records, 10)

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
