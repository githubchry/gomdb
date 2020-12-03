package drivers

import (
	"github.com/minio/minio-go"
	"log"
	"strconv"
)

type MinioCfg struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSL      bool   `json:"ssl"`
}

var MinioDbConn *minio.Client
var MinioDbName string

// 初始化
func MinioDBInit(cfg MinioCfg) error {
	var err error
	// Minio client需要以下4个参数来连接与Amazon S3兼容的对象存储。
	endpoint := cfg.Addr + ":" + strconv.Itoa(cfg.Port) // 对象存储服务的URL
	accessKeyID := cfg.Username                         //Access key是唯一标识你的账户的用户ID。
	secretAccessKey := cfg.Password                     //Secret key是你账户的密码。
	useSSL := cfg.SSL                                   //true代表使用HTTPS

	// 初使化 minio client对象。
	MinioDbConn, err = minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to MinioDB!")

	return err
}

// 关闭
func MinioDBExit() {
	MinioDbConn = nil
	log.Println("MinioDB is closed.")
}
