package create

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/redigo/pkg/prefix"
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

func (r *Redis) Element(key string, val string) error {
	var err error

	var con redis.Conn
	{
		con = r.poo.Get()
		defer con.Close()
	}

	var res string
	{
		res, err = redis.String(con.Do("SET", prefix.WithKeys(r.pre, key), val))
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if res != "OK" {
		return tracer.Mask(executionFailedError)
	}

	return nil
}
