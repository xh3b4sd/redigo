package simple

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo/pkg/prefix"
)

type Search struct {
	pool *redis.Pool

	prefix string
}

func (s *Search) Value(key string) (string, error) {
	con := s.pool.Get()
	defer con.Close()

	res, err := redis.String(con.Do("GET", prefix.WithKeys(s.prefix, key)))
	if IsNotFound(err) {
		return "", tracer.Maskf(notFoundError, key)
	} else if err != nil {
		return "", tracer.Mask(err)
	}

	return res, nil
}
