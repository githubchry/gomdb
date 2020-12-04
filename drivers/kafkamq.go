package drivers

import (
	"github.com/Shopify/sarama"
	"log"
	"strconv"
)

type KafkaCfg struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
}

var KafkaMqClient sarama.Client
var KafkaMqAddr []string

func KafkaMQInit(cfg KafkaCfg) error {

	//详细的config参数:https://blog.csdn.net/chinawangfei/article/details/93097203
	config := sarama.NewConfig()
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	//注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	config.Version = sarama.V2_6_0_0
	config.Net.DialTimeout = 300000000	//3秒 不要用3 * time.Second 对应不上! => 0.3*time.Second, 不设置默认30秒

	KafkaMqAddr = []string{cfg.Addr+":"+strconv.Itoa(cfg.Port)}

	var err error
	log.Println("KafkaMQ Client Conn .....")
	KafkaMqClient, err = sarama.NewClient(KafkaMqAddr, config)
	if err != nil {
		log.Println("create KafkaMq", KafkaMqAddr, "client failed:", err)
		return err
	}

	//获取主题的名称集合
	topics, err := KafkaMqClient.Topics()
	if err != nil {
		log.Println("get topics err:", err)
		return err
	}

	for _, e := range topics {
		log.Println(e)
	}

	//获取broker集合
	brokers := KafkaMqClient.Brokers()
	//输出每个机器的地址
	for _, broker := range brokers {
		log.Println(broker.Addr())
	}

	log.Println("KafkaMQ Init Sucess!")
	//=================================================================================================================
	return err
}

func KafkaMQExit() {

	KafkaMqClient.Close()
}
