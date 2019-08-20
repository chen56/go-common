package redisx

import (
	"github.com/go-redis/redis"
)

// RedisConf ref:github.com/go-redis/redis.Options
type RedisConf struct {
	Addr string `yaml:"addr"          json:"addr"`
	DB   int    `yaml:"db"                json:"db"`
}

type Client struct {
	*redis.Client
}

func (conf RedisConf) NewClient() *Client {
	return &Client{
		redis.NewClient(&redis.Options{
			Addr: conf.Addr,
			DB:   conf.DB,
		}),
	}
}
