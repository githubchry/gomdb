package main

import (
	"bytes"
	"github.com/githubchry/gomdb/drivers"
	"github.com/githubchry/gomdb/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//https://www.cnblogs.com/Dr-wei/p/11742293.html
func main() {
	// log打印设置: Lshortfile文件名+行号  LstdFlags日期加时间
	log.SetFlags(log.Llongfile | log.LstdFlags | log.Lmicroseconds)

	var err error

	//========================================
	minioCfg := drivers.MinioCfg{
		Addr:     "127.0.0.1",
		Port:     9000,		// 默认27017
		Username: "minioadmin",
		Password: "minioadmin",
		SSL:   false,
	}

	// 初始化连接到MongoDB
	err = drivers.MinioDBInit(minioCfg)
	if err != nil {
		log.Fatal(err)
	}

	// 创建存储通
	// 初始化存储桶
	location := "us-east-1"
	bucketNameArr := []string{"image"}

	err = models.MakeBucket(bucketNameArr, location)
	if err != nil {
		log.Fatal(err)
	}

	// 上传本地图片
	models.FPutObject("image", "laopo.jpg" ,"E:\\ubuntu\\codes\\git\\PersonReIDServer\\release\\720.jpg", "image/jpg")

	// 下载图片到本地
	models.FGetObject("image", "laopo.jpg" ,"E:\\ubuntu\\codes\\git\\PersonReIDServer\\release\\laopo.jpg")

	//打开文件句柄操作
	fh, err := os.Open("E:\\ubuntu\\codes\\git\\PersonReIDServer\\release\\laopo.jpg")
	if err != nil {
		log.Fatal("error opening file")
	}
	defer fh.Close()
	img, _ := ioutil.ReadAll(fh)
	fh.Close()

	//上传内存中的文件
	models.PutObject("image", "laopo2.jpg" , "image/jpg", img)

	//下载文件到内存
	img, err = models.GetObject("image", "laopo.jpg" )
	log.Println(err, len(img))

	// 获取下载url
	downloadURL := models.PreDownload("image", "laopo2.jpg")
	log.Println(downloadURL)

	//获取上传url
	uploadURL := models.PreUpload("image", "laopo3.jpg")
	log.Println(uploadURL)

	req, err := http.NewRequest("PUT", uploadURL, bytes.NewBuffer(img))
	if err != nil {
		log.Print(err)
	}

	req.Header.Add("Content-Type","image/jpg")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
	}
	log.Print(resp)

	// 删除单个文件
	err = models.RemoveObject("image", "laopo.jpg")
	if err != nil {
		log.Print(err)
	}

	// 删除多个文件
	objectsCh := make(chan string)
	go func() {
		defer close(objectsCh)
		objectsCh <- "laopo2.jpg"
		objectsCh <- "laopo3.jpg"
	}()
	models.RemoveObjects("image", objectsCh)

	// 断开连接
	drivers.MinioDBExit()
}