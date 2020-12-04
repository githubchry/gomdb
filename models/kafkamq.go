package models

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/githubchry/gomdb/drivers"
	"log"
	"time"
)

func CreateTopics(topics []string) error {

	// 1.连接到Broker, 然后通过这个连接去创建Topic (很扯淡? Topic跟Broker之间不是包含关系, 但只需要理解, Broker只是提交创建Topic的请求到zookeeper, 真正创建的Topic的是zookeeper而不是其下的某一个Broker)
	// [Kafka如何创建topic](https://www.cnblogs.com/warehouse/p/9534230.html)
	// [Kafka解析之topic创建(1)](https://blog.csdn.net/u013256816/article/details/79303825)
	broker, err := drivers.KafkaMqClient.Controller()
	if err != nil {
		log.Fatal(err)
	}

	// check if the connection was OK
	connected, err := broker.Connected()
	if err != nil {
		log.Print(err.Error())
		return err
	}
	log.Print("kafkaMqClient.Controller() connect result:", connected);
	//defer KafkaMqBroker.Close()

	// 2.配置topic参数
	cleanup_policy := "delete"	//自动删除
	retention_bytes := "40000000"	//默认5分钟清空一次大于400M的旧数据
	segment_bytes 	:= "20000000"	//每200M创建新段

	topicDetail := &sarama.TopicDetail{}
	topicDetail.NumPartitions = int32(1)		// 分区数1
	topicDetail.ReplicationFactor = int16(1)	// 复制备份 1
	topicDetail.ConfigEntries = map[string]*string{
		"cleanup.policy" : &cleanup_policy,
		"retention.bytes": &retention_bytes,
		"segment.bytes": &segment_bytes,
	}

	//message_max_bytes := "100000"		// 每条消息最大字节数
	//topicDetail.ConfigEntries["message.max.bytes"] = &message_max_bytes

	topicDetails := make(map[string]*sarama.TopicDetail)
	for _, topic := range topics {
		topicDetails[topic] = topicDetail
	}

	// 创建请求
	request := sarama.CreateTopicsRequest{
		Timeout:      time.Second * 15,
		TopicDetails: topicDetails,
	}

	// 发送请求
	response, err := broker.CreateTopics(&request)
	if err != nil {
		log.Printf("%#v", &err)
		return err
	}

	for key, val := range response.TopicErrors {
		if val.Err != 0 {
			if val.Err == 36 {
				log.Printf("topic [%s] 已存在, 无需重复创建!\n", key)
			} else {
				log.Printf("create topic [%s] error: %s\n", key, val.Error())
				return errors.New(val.Error())
			}
		} else {
			log.Printf("topic [%s] 创建成功!\n", key)
		}
	}

	return nil
}

func DeleteTopics(topics []string) {
	// 创建删除请求
	request := sarama.DeleteTopicsRequest{
		Timeout:	time.Second * 15,
		Topics: topics,
	}

	broker, err := drivers.KafkaMqClient.Controller()
	if err != nil {
		log.Println(err)
		return
	}

	response, err := broker.DeleteTopics(&request)
	if err != nil {
		log.Printf("%#v", &err)
	}

	for key, val := range response.TopicErrorCodes {
		if val != 0 {
			log.Printf("delete topic [%s] error: %s\n", key, val.Error())
		} else {
			log.Printf("topic [%s] 已删除!\n", key)
		}
	}
}



func LoopConsumer(consumer sarama.Consumer, topic string, process func(msg *sarama.ConsumerMessage)) {
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		log.Printf("fail to get list of %v partition:%v\n", topic, err)
	}

	// partition号从0开始
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		partitionConsumer, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			log.Printf("failed to start consumer for %v partition %d,err:%v\n", topic, partition, err)
			return
		}
		defer partitionConsumer.AsyncClose()

		// 异步从每个分区消费信息
		func(sarama.PartitionConsumer) {
			for msg := range partitionConsumer.Messages() {
				//log.Printf("Partition:%d Offset:%d Key:%v len(Value):%v", msg.Partition, msg.Offset, msg.Key, len(msg.Value))
				process(msg)
			}
		}(partitionConsumer)
	}
}

func CreateConsumer(config sarama.Config) (sarama.Consumer, error) {
	consumer, err := sarama.NewConsumer(drivers.KafkaMqAddr, &config)
	if err != nil {
		log.Println("NewSyncProducer", err)
	}
	return consumer, err
}

func CreateProducer(config sarama.Config) (sarama.AsyncProducer, error) {
	//使用配置,新建一个异步生产者
	producer, err := sarama.NewAsyncProducer(drivers.KafkaMqAddr, &config)
	if err != nil {
		log.Println("NewSyncProducer", err)
	}

	return producer, err
}

func ProducerInput(producer sarama.AsyncProducer, msg *sarama.ProducerMessage) (int64, error) {
	//使用通道发送
	producer.Input() <- msg
	//循环判断哪个通道发送过来数据.
	select {
	case sucess := <-producer.Successes():
		//log.Println("offset: ", sucess.Offset, "timestamp: ", sucess.Timestamp.String(), "partitions: ", sucess.Partition)
		return sucess.Offset, nil
	case fail := <-producer.Errors():
		log.Println("err: ", fail.Err)
		return -1, fail.Err
	}
}