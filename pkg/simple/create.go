package simple

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo/pkg/prefix"
)

type create struct {
	pool *redis.Pool

	prefix string
}

func (c *create) Element(key string, val string) error {
	con := c.pool.Get()
	defer con.Close()

	reply, err := redis.String(con.Do("SET", prefix.WithKeys(c.prefix, key), val))
	if err != nil {
		return tracer.Mask(err)
	}

	if reply != "OK" {
		return tracer.Mask(executionFailedError)
	}

	return nil
}
