package kafka

import (
	"testing"
	"time"

	"github.com/chen56/go-common/assert"
	"github.com/Shopify/sarama"
	"github.com/apex/log"
)

func TestManual_NewBatchConsumer(t *testing.T) {
	config := NewConfig()
	config.Consumer.Offsets.Retention = 5 * time.Minute
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true
	//config.Consumer.MaxWaitTime=3*time.Second
	config.Consumer.Fetch.Min = 100000
	//config.Consumer.Fetch.Default = 512000

	config.Group.Return.Notifications = true
	config.MarkOffset = false
	config.BatchLimit=100
	config.BatchFetchTimeout=1*time.Second
	consumer, err := NewBatchConsumer([]string{"192.168.1.11:9092"}, "chenpeng2", []string{"appshare-db"}, config, func(msgs []*sarama.ConsumerMessage) error {
		log.Infof("len: %+v", len(msgs))
		for _, msgRaw := range msgs {
			//msgUnmarshal:=config.Match(string(msgRaw.Value))
			log.Infof("msgRaw: %+v", string(msgRaw.Value))
		}
		return nil
	})
	defer consumer.Close()
	assert.NoErr(err)
	consumer.Run()
}
