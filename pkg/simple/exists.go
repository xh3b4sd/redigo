package simple

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo/pkg/prefix"
)

type Exists struct {
	pool *redis.Pool

	prefix string
}

func (e *Exists) Element(key string) (bool, error) {
	con := e.pool.Get()
	defer con.Close()

	result, err := redis.Bool(con.Do("EXISTS", prefix.WithKeys(e.prefix, key)))
	if err != nil {
		return false, tracer.Mask(err)
	}

	return result, nil
}
