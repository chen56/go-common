package main

import (
	"fmt"
	"time"
	"github.com/chen56/go-common/must"
	"bufio"
	"os"

	"github.com/chen56/go-common/kafka"
	"github.com/Shopify/sarama"
)
type empty struct{}


func main() {
	config := kafka.NewConfig()
	config.Consumer.Offsets.Retention = 5 * time.Minute
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.MarkOffset = false

	consumer, err := kafka.NewBatchConsumer([]string{"192.168.1.11:9092"}, "chenpeng2", []string{"x"}, config, func(msgs []*sarama.ConsumerMessage) error {
		fmt.Printf("len: %+v", len(msgs))
		for _, msgRaw := range msgs {
			//msgUnmarshal:=config.Match(string(msgRaw.Value))
			fmt.Printf("msgRaw: %+v\n", string(msgRaw.Value))
		}
		return nil
	})
	defer consumer.Close()
	must.NoErr(err)

	//暂停等待输入
	fmt.Printf("waiting...")
	reader := bufio.NewReader(os.Stdin)
	strBytes, hasMore, err := reader.ReadLine()
	fmt.Printf("run consumter %v %v %v", strBytes, hasMore, err)

	//期望接收到等待时发出的消息
	consumer.Run()
}
