package kafka

import (
	"github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
	"github.com/apex/log"
	"github.com/Shopify/sarama"
)

type Callback func(msgs []*sarama.ConsumerMessage) error;

type BatchConsumer struct {
	consumer *cluster.Consumer
	config   *BatchConfig
	callback Callback
	addrs    []string
	groupID  string
	topics   []string
}
type BatchConfig struct {
	cluster.Config
	BatchLimit int
	MarkOffset bool
}

func NewConfig() *BatchConfig {
	c := &BatchConfig{
		Config: *cluster.NewConfig(),
	}
	c.BatchLimit = 1
	return c
}

func NewBatchConsumer(addrs []string, groupID string, topics []string, config *BatchConfig, callback Callback) (*BatchConsumer, error) {
	if callback == nil {
		panic("Callback should not be nil")
	}
	if config.BatchLimit <= 0 {
		panic("BatchLimit should >0")
	}

	consumer, err := cluster.NewConsumer(addrs, groupID, topics, &config.Config)
	if err != nil {
		return nil, err
	}
	return &BatchConsumer{
		consumer: consumer,
		config:   config,
		callback: callback,
		addrs:    addrs,
		groupID:  groupID,
		topics:   topics,
	}, nil
}

func (this *BatchConsumer) Close() error {
	return this.consumer.Close()
}

func (this *BatchConsumer) Run() {
	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var log = log.WithField("addrs", this.addrs).
		WithField("class", "BatchConsumer").
		WithField("groupID", this.groupID).
		WithField("topics", this.topics)

	log.Info("Run")
	// consume errors
	go func() {
		for err := range this.consumer.Errors() {
			log.Errorf("TickError: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range this.consumer.Notifications() {
			log.Infof("Rebalanced: %+v\n", ntf)
		}
	}()

	var msgs []*sarama.ConsumerMessage
	var i int
	for {
		select {
		case msg, ok := <-this.consumer.Messages():
			if !ok {
				continue
			}
			log.Infof("received: %s", string(msg.Value))

			msgs = append(msgs, msg)

			i++
			if i%this.config.BatchLimit != 0 {
				continue
			}

			err := this.callback(msgs)
			if err != nil {
				log.Errorf("batch consumer error: %+v", err)
				return
			}

			msgs = []*sarama.ConsumerMessage{}
			if this.config.MarkOffset {
				this.consumer.MarkOffset(msg, "") // mark message as processed
			}
		case s, ok := <-signals:
			if ok {
				log.Infof("os.Signal: %+v", s)
			} else {
				log.Info("os.Signal: what s happen?")
			}
			return
		}
	}
}
