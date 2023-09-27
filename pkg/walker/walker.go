package walker

import (
	"github.com/gomodule/redigo/redis"
)

type Config struct {
	Pool *redis.Pool

	Count  int
	Prefix string
}

type Walker struct {
	pool *redis.Pool

	count  int
	prefix string
}

func New(config Config) (*Walker, error) {
	if config.Count <= 0 {
		config.Count = 100
	}

	w := &Walker{
		pool: config.Pool,

		count:  config.Count,
		prefix: config.Prefix,
	}

	return w, nil
}
