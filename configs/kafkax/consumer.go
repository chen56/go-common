package kafkax

import (
	"time"
	"context"
	"strings"
	//"runtime"
	"sync"
	"syscall"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

//消费者配置
type ConsumerConf struct {
	Brokers      string `yaml:"brokers"      json:"brokers"`
	Version      string `yaml:"version"      json:"version"`
	Group        string `yaml:"group"        json:"group"`
	ClientID     string `yaml:"clientID"     json:"clientID"`
	Oldest       bool   `yaml:"oldest"       json:"oldest"`
	MaxWaitTimeMs    int64  `yaml:"maxWaitTimeMs"    json:"maxWaitTimeMs"` //fetch.wait.max.ms
	SessionTimeoutMs int64  `yaml:"sessionTimeoutMs" json:"sessionTimeoutMs"` //session.timeout.ms
	KeepAliveMs  int64  `yaml:"keepAliveMs"  json:"keepAliveMs"`

	Logger sarama.StdLogger
}

func NewConsumerConf() *ConsumerConf {
	return &ConsumerConf{
		Brokers:       "localhost:9092",
		Version:       "2.1.1",
		Group:         "chen56-consumer-group",
		ClientID:      "chen56-consumer-client",
		Oldest:        false,
		MaxWaitTimeMs:    2000,
		SessionTimeoutMs: 300000,
		KeepAliveMs:      86400000,
	}
}

//创建消费者
func (x ConsumerConf) NewConsumer() (*Consumer) {
	if x.Logger != nil {
		sarama.Logger = x.Logger
	}

	version, err := sarama.ParseKafkaVersion(x.Version)
	if err != nil {
		panic(err)
	}

	config := sarama.NewConfig()
	config.Version = version
	config.ClientID = x.ClientID
	if x.Oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	config.Consumer.MaxWaitTime = time.Duration(x.MaxWaitTimeMs) * time.Millisecond
	config.Consumer.Group.Session.Timeout = time.Duration(x.SessionTimeoutMs) * time.Millisecond
	config.Net.KeepAlive = time.Duration(x.KeepAliveMs) * time.Millisecond

	client, err := sarama.NewConsumerGroup(strings.Split(x.Brokers, ","), x.Group, config)
	if err != nil {
		panic(err)
	}
	//runtime.SetFinalizer(client, func(client *sarama.ConsumerGroup) {
	//	client.Close()
	//})

	return &Consumer{
		ConsumerGroup: client,
		ready:  make(chan bool),
		conf:   &x,
	}
}

// 消费者
type Consumer struct {
	sarama.ConsumerGroup
	ready    chan bool
	conf     *ConsumerConf
	callback func(*ConsumerMessage)(bool)
}

type ConsumerMessage struct {
	*sarama.ConsumerMessage 
}

//开始消费
func (consumer *Consumer) Consume(topics string, callback func(*ConsumerMessage)(bool)) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	consumer.callback = callback

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := consumer.ConsumerGroup.Consume(ctx, strings.Split(topics, ","), consumer); err != nil {
				sarama.Logger.Println(err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	sarama.Logger.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		sarama.Logger.Println("terminating: context cancelled")
	case <-sigterm:
		sarama.Logger.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err := consumer.Close(); err != nil {
		return err
	}
	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		//sarama.Logger.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		success := consumer.callback(&ConsumerMessage{message})
		if success {
			session.MarkMessage(message, "")
		}
	}

	return nil
}
