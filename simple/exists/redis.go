package exists

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

func (r *Redis) Multi(key ...string) (int64, error) {
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

	var res int64
	{
		res, err = redis.Int64(con.Do("EXISTS", arg...))
		if err != nil {
			return 0, tracer.Mask(err)
		}
	}

	return res, nil
}
