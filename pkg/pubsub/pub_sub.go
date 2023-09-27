package pubsub

import (
	"github.com/gomodule/redigo/redis"
)

type Config struct {
	Pool *redis.Pool

	Prefix string
}

type PubSub struct {
	pool *redis.Pool

	prefix string
}

func New(config Config) (*PubSub, error) {
	p := &PubSub{
		pool: config.Pool,

		prefix: config.Prefix,
	}

	return p, nil
}
