package main

import (
	"github.com/Shopify/sarama"
	"github.com/githubchry/gomdb/drivers"
	"github.com/githubchry/gomdb/models"
	"log"
)

//https://www.cnblogs.com/Dr-wei/p/11742293.html
func main() {
	// log打印设置: Lshortfile文件名+行号  LstdFlags日期加时间
	log.SetFlags(log.Llongfile | log.LstdFlags | log.Lmicroseconds)

	var err error

	//========================================
	kafkaCfg := drivers.KafkaCfg{
		Addr:     "127.0.0.1",
		Port:     9092,
	}

	// 初始化连接到MongoDB
	err = drivers.KafkaMQInit(kafkaCfg)
	if err != nil {
		log.Fatal(err)
	}

	//创建 topic
	models.CreateTopics([]string{"hello", "world", "test"})

	//创建消息生产者
	//详细的config参数:https://blog.csdn.net/chinawangfei/article/details/93097203
	config := sarama.NewConfig()
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	//注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	config.Version = sarama.V2_6_0_0
	//============ Producer config ============
	//随机向partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//等待服务器所有副本都保存成功后的响应 	发送完数据需要leader和follow都确认
	config.Producer.RequiredAcks = sarama.WaitForAll
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse(默认)这里才有用. 成功交付的消息将在success channel返回
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//============ Consumer config ============
	config.Consumer.Return.Errors = true
	//config.Consumer.Fetch.Max = 16
	//config.Consumer.Fetch.Min = 16
	//config.Consumer.MaxWaitTime = time.Second * 10

	testProducer, err := models.CreateProducer(*config)
	if err != nil {
		log.Println("Create Consumer", err)
	}

	// 创建消费者
	testConsumer, err := models.CreateConsumer(*config)
	if err != nil {
		log.Println("Create Image Consumer", err)
	}

	//消费消息线程
	go models.LoopConsumer(testConsumer, "test", func(msg *sarama.ConsumerMessage) {
		log.Printf("从kafka收到消息[%s], offset[%d]\n", string(msg.Value), msg.Offset)
	})

	// 生产消息
	msg := &sarama.ProducerMessage{
		Topic : "test",
		Value : sarama.StringEncoder("hello world"),
	}
	_, err = models.ProducerInput(testProducer, msg)
	if err != nil {
		log.Println("Producer failed:", err)
	} else {
		log.Println("已经发送消息到kafka...")
	}


	//删除 topic
	models.DeleteTopics([]string{"hello", "world", "test"})

	// 断开连接
	drivers.KafkaMQExit()
}