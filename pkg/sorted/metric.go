package sorted

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/redigo/prefix"
	"github.com/xh3b4sd/tracer"
)

type metric struct {
	pool *redis.Pool

	prefix string
}

func (m *metric) Count(key string) (int64, error) {
	con := m.pool.Get()
	defer con.Close()

	res, err := redis.Int64(con.Do("ZCARD", prefix.WithKeys(m.prefix, key)))
	if err != nil {
		return 0, tracer.Mask(err)
	}

	return res, nil
}
