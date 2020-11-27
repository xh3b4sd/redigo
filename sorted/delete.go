package sorted

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

type Delete struct {
	pool *redis.Pool

	prefix string
}

func (d *Delete) Value(key string, val string) error {
	con := d.pool.Get()
	defer con.Close()

	_, err := redis.Int(con.Do("ZREM", withPrefix(d.prefix, key), val))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
