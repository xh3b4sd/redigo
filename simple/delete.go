package simple

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo/prefix"
)

type Delete struct {
	pool *redis.Pool

	prefix string
}

func (d *Delete) Element(key string) error {
	con := d.pool.Get()
	defer con.Close()

	_, err := redis.Int64(con.Do("DEL", prefix.WithKeys(d.prefix, key)))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
