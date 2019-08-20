package fluentx

import (
	"github.com/chen56/go-common/must"
	"github.com/fluent/fluent-logger-golang/fluent"
	"time"
)

// copy from: fluent.Config
type FluentConf struct {
	FluentPort         int           `yaml:"fluent_port"          json:"fluent_port"`
	FluentHost         string        `yaml:"fluent_host"          json:"fluent_host"`
	FluentNetwork      string        `yaml:"fluent_network"       json:"fluent_network"`
	FluentSocketPath   string        `yaml:"fluent_socket_path"   json:"fluent_socket_path"`
	Timeout            time.Duration `yaml:"timeout"              json:"timeout"`
	WriteTimeout       time.Duration `yaml:"write_timeout"        json:"write_timeout"`
	BufferLimit        int           `yaml:"buffer_limit"         json:"buffer_limit"`
	RetryWait          int           `yaml:"retry_wait"           json:"retry_wait"`
	MaxRetry           int           `yaml:"max_retry"            json:"max_retry"`
	MaxRetryWait       int           `yaml:"max_retry_wait"       json:"max_retry_wait"`
	TagPrefix          string        `yaml:"tag_prefix"           json:"tag_prefix"`
	Async              bool          `yaml:"async"                json:"async"`
	MarshalAsJSON      bool          `yaml:"marshal_as_json"      json:"marshal_as_json"`
	SubSecondPrecision bool          `yaml:"sub_second_precision" json:"sub_second_precision"`
	RequestAck         bool          `yaml:"request_ack"          json:"request_ack"`
}

func (x FluentConf) MustNewFluent() *fluent.Fluent {
	fConfig := fluent.Config{
		FluentHost:         x.FluentHost,
		FluentPort:         x.FluentPort,
		FluentNetwork:      x.FluentNetwork,
		FluentSocketPath:   x.FluentSocketPath,
		Timeout:            x.Timeout,
		WriteTimeout:       x.WriteTimeout,
		BufferLimit:        x.BufferLimit,
		RetryWait:          x.RetryWait,
		MaxRetry:           x.MaxRetry,
		MaxRetryWait:       x.MaxRetryWait,
		TagPrefix:          x.TagPrefix,
		Async:              x.Async,
		MarshalAsJSON:      x.MarshalAsJSON,
		SubSecondPrecision: x.SubSecondPrecision,
		RequestAck:         x.RequestAck,
	}
	logger, err := fluent.New(fConfig)
	must.NoError(err)
	return logger
}
