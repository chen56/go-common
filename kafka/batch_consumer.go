package kafka

import (
	"github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
	"github.com/apex/log"
	"github.com/Shopify/sarama"
	"reflect"
)

type empty struct{}
var ll *log.Entry = log.WithField("pkg", reflect.TypeOf(empty{}).PkgPath())

type Callback func(msgs []*sarama.ConsumerMessage) error;

type BatchConsumer struct {
	consumer *cluster.Consumer
	config   *BatchConfig
	callback Callback
	addrs    []string
	groupID  string
	topics   []string
	ll *log.Entry
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

	var ll = ll.WithField("addrs", addrs).
		WithField("groupID", groupID).
		WithField("topics", topics)

	consumer, err := cluster.NewConsumer(addrs, groupID, topics, &config.Config)
	if err != nil {
		return nil, err
	}

	//原先consumer.Errors和consumer.Notifications放在Run(),但发现不能保证delay_read
	// consume errors
	go func() {
		for err := range consumer.Errors() {
			ll.Errorf("TickError: %s\n", err.Error())
		}
	}()
	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			ll.Infof("Rebalanced: %+v\n", ntf)
		}
	}()

	return &BatchConsumer{
		consumer: consumer,
		config:   config,
		callback: callback,
		addrs:    addrs,
		groupID:  groupID,
		topics:   topics,
		ll:ll,
	},nil
}

func (this *BatchConsumer) Close() error {
	return this.consumer.Close()
}

func (x *BatchConsumer) Run() {
	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)


	var msgs []*sarama.ConsumerMessage
	var i int
	for {
		select {
		case msg, ok := <-x.consumer.Messages():
			if !ok {
				continue
			}
		        if x.ll.Level==log.DebugLevel{
				x.ll.Debugf("received: %s", string(msg.Value))
			}

			msgs = append(msgs, msg)

			i++
			if i% x.config.BatchLimit != 0 {
				continue
			}

			err := x.callback(msgs)
			if err != nil {
				x.ll.Errorf("batch consumer error: %+v", err)
				return
			}

			msgs = []*sarama.ConsumerMessage{}
			if x.config.MarkOffset {
				x.consumer.MarkOffset(msg, "") // mark message as processed
			}
		case s, ok := <-signals:
			if ok {
				x.ll.Infof("os.Signal: %+v", s)
			} else {
				x.ll.Info("os.Signal: what s happen?")
			}
			return
		}
	}
}
