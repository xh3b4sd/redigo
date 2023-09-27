package sorted

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo/pkg/index"
	"github.com/xh3b4sd/redigo/pkg/prefix"
)

type exists struct {
	pool *redis.Pool

	prefix string
}

func (e *exists) Index(key string, ind string) (bool, error) {
	con := e.pool.Get()
	defer con.Close()

	_, err := redis.Bytes(con.Do("ZSCORE", prefix.WithKeys(e.prefix, index.New(key)), ind))
	if IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, tracer.Mask(err)
	}

	return true, nil
}

func (e *exists) Score(key string, sco float64) (bool, error) {
	con := e.pool.Get()
	defer con.Close()

	res, err := redis.Strings(con.Do("ZRANGE", prefix.WithKeys(e.prefix, key), sco, sco, "BYSCORE"))
	if err != nil {
		return false, tracer.Mask(err)
	}
	if len(res) == 0 {
		return false, nil
	}

	return true, nil
}

func (e *exists) Value(key string, val string) (bool, error) {
	con := e.pool.Get()
	defer con.Close()

	_, err := redis.Bytes(con.Do("ZSCORE", prefix.WithKeys(e.prefix, key), val))
	if IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, tracer.Mask(err)
	}

	return true, nil
}
