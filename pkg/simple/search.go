package simple

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo/pkg/prefix"
)

type search struct {
	pool *redis.Pool

	prefix string
}

func (s *search) Multi(key ...string) ([]string, error) {
	con := s.pool.Get()
	defer con.Close()

	var mul []interface{}
	for _, x := range key {
		mul = append(mul, prefix.WithKeys(s.prefix, x))
	}

	res, err := redis.Strings(con.Do("MGET", mul...))
	if err != nil {
		return nil, tracer.Mask(err)
	}

	if len(res) == 1 && res[0] == "" {
		return nil, tracer.Mask(notFoundError)
	}

	return res, nil
}
