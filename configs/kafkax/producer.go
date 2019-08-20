package kafkax

import (
	"time"
	"strings"
	//"runtime"

	"github.com/Shopify/sarama"
)

//生产者配置
type ProducerConf struct {
	Hosts        string `yaml:"hosts"        json:"hosts"`
	ClientID     string `yaml:"clientID"     json:"clientID"`
	KeepAliveMs  int64  `yaml:"keepAliveMs"  json:"keepAliveMs"`
	ReqTimeoutMs int64  `yaml:"reqTimeoutMs" json:"reqTimeoutMs"`
	RetSuccesses bool   `yaml:"retSuccesses" json:"retSuccesses"`
}

//创建生产者配置
func NewProducerConf() *ProducerConf {
	return &ProducerConf{
		Hosts:        "localhost:9092",
		ClientID:     "chen56-producer-client",
		KeepAliveMs:  86400000,
		ReqTimeoutMs: 300000,
		RetSuccesses: false,
	}
}

//创建Kafka客户端
func (x ProducerConf) NewClient() (client sarama.Client) {
	config := sarama.NewConfig()
	config.ClientID = x.ClientID
	config.Net.KeepAlive = time.Duration(x.KeepAliveMs) * time.Millisecond
	config.Producer.Timeout = time.Duration(x.ReqTimeoutMs) * time.Millisecond
	config.Producer.Return.Successes = x.RetSuccesses
	client, err := sarama.NewClient(strings.Split(x.Hosts, ","), config)
	if err != nil {
		panic(err)
	}
	return
}

//创建异步生产者
func (x ProducerConf) NewAsyncProducer() *AsyncProducer {
	client := x.NewClient()
	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	//runtime.SetFinalizer(producer, func(producer *sarama.AsyncProducer) {
	//	producer.Close()
	//})
	return &AsyncProducer{producer}
}

//创建同步生产者
func (x ProducerConf) NewSyncProducer() *SyncProducer {
	x.RetSuccesses = true //必须为true
	client := x.NewClient()
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	//runtime.SetFinalizer(producer, func(producer *sarama.SyncProducer) {
	//	producer.Close()
	//})
	return &SyncProducer{producer}
}

//异步生产者
type AsyncProducer struct {
	sarama.AsyncProducer
}

//异步发送消息
func (producer AsyncProducer) SendMessage(topic, text string) {
	producer.AsyncProducer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder(text)}
}

//异步发送带Key消息
//func (producer AsyncProducer) SendMessageWithKey(topic, key, text string) {
//	producer.AsyncProducer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: key, Value: sarama.StringEncoder(text)}
//}

//同步生产者
type SyncProducer struct {
	sarama.SyncProducer
}

//同步发送消息
func (producer SyncProducer) SendMessage(topic, text string) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key: nil,
		Value: sarama.StringEncoder(text),
	}
	return producer.SyncProducer.SendMessage(msg)
}

//同步发送带Key消息
//func (producer SyncProducer) SendMessageWithKey(topic, key, text string) (partition int32, offset int64, err error) {
//	msg := &sarama.ProducerMessage{
//		Topic: topic,
//		Key: key,
//		Value: sarama.StringEncoder(text)
//	}
//	return producer.SyncProducer.SendMessage(msg)
//}
