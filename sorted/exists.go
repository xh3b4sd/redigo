package sorted

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

type Exists struct {
	pool *redis.Pool

	prefix string
}

func (e *Exists) Score(key string, sco float64) (bool, error) {
	con := e.pool.Get()
	defer con.Close()

	res, err := redis.Strings(con.Do("ZRANGEBYSCORE", withPrefix(e.prefix, key), sco, sco))
	if err != nil {
		return false, tracer.Mask(err)
	}
	if len(res) == 0 {
		return false, nil
	}

	return true, nil
}
