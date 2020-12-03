package drivers

import (
	"github.com/milvus-io/milvus-sdk-go/milvus"
	"log"
	"strconv"
)

type MilvusCfg struct {
	Addr     string `json:"addr"`
	Port     int 	`json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}


var MilvusDbConn milvus.MilvusClient

// 初始化
func MilvusDBInit(cfg MilvusCfg) error {
	var grpcClient milvus.Milvusclient
	MilvusDbConn = milvus.NewMilvusClient(grpcClient.Instance)

	//Client version
	println("Client version: " + MilvusDbConn.GetClientVersion())

	//test connect
	connectParam := milvus.ConnectParam{cfg.Addr, strconv.Itoa(cfg.Port)}
	err := MilvusDbConn.Connect(connectParam)
	if err != nil {
		println("client: connect failed: " + err.Error())
	}

	if MilvusDbConn.IsConnected() == false {
		println("client: not connected: ")
		return err
	}
	println("Server status: connected")
	return err
}

// 关闭
func MilvusDBExit() {
	err := MilvusDbConn.Disconnect()
	if err != nil {
		println("Disconnect failed!")
		return
	}
	println("Client disconnect server success!")
	log.Println("MilvusDB is closed.")
}


