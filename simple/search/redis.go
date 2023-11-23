package search

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/redigo/prefix"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Poo *redis.Pool
	Pre string
}

type Redis struct {
	poo *redis.Pool
	pre string
}

func New(c Config) *Redis {
	return &Redis{
		poo: c.Poo,
		pre: c.Pre,
	}
}

func (r *Redis) Multi(key ...string) ([]string, error) {
	var err error

	var con redis.Conn
	{
		con = r.poo.Get()
		defer con.Close()
	}

	var arg []interface{}
	for _, x := range key {
		arg = append(arg, prefix.WithKeys(r.pre, x))
	}

	var res []string
	{
		res, err = redis.Strings(con.Do("MGET", arg...))
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	if len(res) == 1 && res[0] == "" {
		return nil, tracer.Mask(notFoundError)
	}

	return res, nil
}
