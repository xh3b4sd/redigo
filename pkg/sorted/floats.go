package sorted

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo/pkg/prefix"
)

type floats struct {
	pool *redis.Pool

	prefix string
}

func (f *floats) Score(key string, val string, sco float64) (float64, error) {
	con := f.pool.Get()
	defer con.Close()

	res, err := redis.Float64(con.Do("ZINCRBY", prefix.WithKeys(f.prefix, key), sco, val))
	if err != nil {
		return 0, tracer.Mask(err)
	}

	return res, nil
}
