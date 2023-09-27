package simple

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo/pkg/prefix"
)

type delete struct {
	pool *redis.Pool

	prefix string
}

func (d *delete) Multi(key ...string) (int64, error) {
	con := d.pool.Get()
	defer con.Close()

	var mul []interface{}
	for _, x := range key {
		mul = append(mul, prefix.WithKeys(d.prefix, x))
	}

	res, err := redis.Int64(con.Do("DEL", mul...))
	if err != nil {
		return 0, tracer.Mask(err)
	}

	return res, nil
}
