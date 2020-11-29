package sorted

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo/prefix"
)

type Exists struct {
	pool *redis.Pool

	prefix string
}

func (e *Exists) Score(key string, sco float64) (bool, error) {
	con := e.pool.Get()
	defer con.Close()

	res, err := redis.Strings(con.Do("ZRANGEBYSCORE", prefix.WithKeys(e.prefix, key), sco, sco))
	if err != nil {
		return false, tracer.Mask(err)
	}
	if len(res) == 0 {
		return false, nil
	}

	return true, nil
}

func (e *Exists) Value(key string, val string) (bool, error) {
	con := e.pool.Get()
	defer con.Close()

	_, err := redis.Strings(con.Do("ZSCORE", prefix.WithKeys(e.prefix, key)))
	if IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, tracer.Mask(err)
	}

	return true, nil
}
