package kafka

import (
	"testing"
	"time"

	"github.com/chen56/go-common/assert"
	"github.com/Shopify/sarama"
	"github.com/apex/log"
	"github.com/chen56/go-common/testx"
)

func TestManual_NewBatchConsumer(t *testing.T) {
	testx.Skip(t)
	config := NewConfig()
	config.Consumer.Offsets.Retention = 5 * time.Minute
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.MarkOffset = false

	consumer, err := NewBatchConsumer([]string{"192.168.1.11:9092"}, "chenpeng", []string{"member-event"}, config, func(msgs []*sarama.ConsumerMessage) error {
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
