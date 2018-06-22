package kafka

import (
	"github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
	"github.com/apex/log"
	"github.com/Shopify/sarama"
	"reflect"
	"time"
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
	ll       *log.Entry
}
type BatchConfig struct {
	cluster.Config
	MarkOffset bool
	//等待下批数据到达BatchLimit上限后再返回这批数据
	BatchLimit int
	//等待下批数据到达BatchLimit上限前，如果超时，也返回
	BatchFetchTimeout time.Duration
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
	}, nil
}

func (this *BatchConsumer) Close() error {
	return this.consumer.Close()
}

func (x *BatchConsumer) Run()error {
	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var bufferMessages []*sarama.ConsumerMessage

	var i int
	for {
		select {
		case msg, ok := <-x.consumer.Messages():
			if !ok {
				continue
			}
			if x.ll.Level == log.DebugLevel {
				x.ll.Debugf("received: %s", string(msg.Value))
			}

			bufferMessages = append(bufferMessages, msg)

			i++
			if i % x.config.BatchLimit != 0 {
				continue
			}

			err:=x.process(&bufferMessages)
			if err != nil {
				return err
			}
		case s, ok := <-signals:
			if ok {
				x.ll.Infof("os.Signal: %+v", s)
			} else {
				x.ll.Info("os.Signal: what s happen?")
			}
			return nil
		case <-time.After(3 * time.Second):
			err:=x.process(&bufferMessages)
			if err != nil {
				return err
			}
		}

	}
}

func (x *BatchConsumer)  process(bufferMessages *[]*sarama.ConsumerMessage)error{
	err := x.callback(*bufferMessages)
	if err != nil {
		return err
	}
	*bufferMessages = []*sarama.ConsumerMessage{}
	if x.config.MarkOffset {
		for _,msg:=range *bufferMessages {
			x.consumer.MarkOffset(msg, "") // mark message as processed
		}
	}
	return nil
}