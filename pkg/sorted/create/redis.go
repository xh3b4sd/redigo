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
		for _, x := range ind {
			m[x] = m[x] + 1
		}

		for _, v := range m {
			if v > 1 {
				return tracer.Maskf(executionFailedError, "index must be unique")
			}
		}

		for _, x := range ind {
			if x == "" {
				return tracer.Maskf(executionFailedError, "index must not be empty")
			}
			if strings.Count(x, " ") != 0 {
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

		for _, x := range ind {
			arg = append(arg, x)
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

	// Use NX to only add new elements. Don't update already existing elements.
	var res int64
	{
		res, err = redis.Int64(con.Do("ZADD", prefix.WithKeys(r.pre, key), "NX", sco, val))
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// ZADD returns the number of elements created. Scores can be duplicated, but
	// values must be unique after all. If the value to be added existed already,
	// then ZADD returns 0. In that case we return the appropriate error.
	if res == 0 {
		return tracer.Maskf(alreadyExistsError, "key %s does already hold a value for score %f", key, sco)
	}

	return nil
}

func (r *Redis) Union(dst string, key ...string) (int64, error) {
	var err error

	var con redis.Conn
	{
		con = r.poo.Get()
		defer con.Close()
	}

	var arg []interface{}
	{
		arg = append(arg, dst, len(key))

		for _, x := range key {
			arg = append(arg, prefix.WithKeys(r.pre, x))
		}

		arg = append(arg, "AGGREGATE", "MIN")
	}

	var res int64
	{
		res, err = redis.Int64(con.Do("ZUNIONSTORE", arg...))
		if err != nil {
			return 0, tracer.Mask(err)
		}
	}

	return res, nil
}
