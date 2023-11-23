package create

import (
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/redigo/index"
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
	scr *redis.Script
}

func New(c Config) *Redis {
	return &Redis{
		poo: c.Poo,
		pre: c.Pre,
		scr: redis.NewScript(2, createIndexScript),
	}
}

func (r *Redis) Index(key string, val string, sco float64, ind ...string) error {
	var err error

	var con redis.Conn
	{
		con = r.poo.Get()
		defer con.Close()
	}

	if len(ind) != 0 {
		m := map[string]int{}
		for _, s := range ind {
			m[s] = m[s] + 1
		}

		for _, v := range m {
			if v > 1 {
				return tracer.Maskf(executionFailedError, "index must be unique")
			}
		}

		for _, s := range ind {
			if s == "" {
				return tracer.Maskf(executionFailedError, "index must not be empty")
			}
			if strings.Count(s, " ") != 0 {
				return tracer.Maskf(executionFailedError, "index must not contain whitespace")
			}
		}
	}

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(r.pre, key))            // KEYS[1]
		arg = append(arg, prefix.WithKeys(r.pre, index.New(key))) // KEYS[2]
		arg = append(arg, val)                                    // ARGV[1]
		arg = append(arg, sco)                                    // ARGV[2]

		for _, s := range ind {
			arg = append(arg, s)
		}
	}

	var res int64
	{
		res, err = redis.Int64(r.scr.Do(con, arg...))
		if err != nil {
			return tracer.Mask(err)
		}
	}

	switch res {
	case 0:
		return tracer.Maskf(alreadyExistsError, "score must be unique")
	case 1:
		return tracer.Maskf(alreadyExistsError, "index must be unique")
	case 2:
		return nil
	}

	return tracer.Mask(executionFailedError)
}

func (r *Redis) Score(key string, val string, sco float64) error {
	var err error

	var con redis.Conn
	{
		con = r.poo.Get()
		defer con.Close()
	}

	var res int64
	{
		res, err = redis.Int64(con.Do("ZADD", prefix.WithKeys(r.pre, key), "NX", sco, val))
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// ZADD returns the number of elements created. Scores can be duplicated, but
	// values must be unique after all. If the value to be added existed already,
	// ZADD returns 0. In that case we return the appropriate error.
	if res == 0 {
		return tracer.Maskf(alreadyExistsError, "value must be unique")
	}

	return nil
}
