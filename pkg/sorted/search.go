package sorted

import (
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo/pkg/index"
	"github.com/xh3b4sd/redigo/pkg/prefix"
)

const indexElementScript = `
	local sco = redis.call("ZSCORE", KEYS[2], ARGV[1])
	local res = redis.call("ZRANGEBYSCORE", KEYS[1], sco, sco)

	return res[1]
`

type Search struct {
	pool *redis.Pool

	indexElementScript *redis.Script

	prefix string
}

func (s *Search) Index(key string, ind string) (string, error) {
	con := s.pool.Get()
	defer con.Close()

	{
		if ind == "" {
			return "", tracer.Maskf(executionFailedError, "index must not be empty")
		}
		if strings.Count(ind, " ") != 0 {
			return "", tracer.Maskf(executionFailedError, "index must not contain whitespace")
		}
	}

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(s.prefix, key))
		arg = append(arg, prefix.WithKeys(s.prefix, index.New(key)))
		arg = append(arg, ind)
	}

	res, err := redis.String(s.indexElementScript.Do(con, arg...))
	if err != nil {
		return "", tracer.Mask(err)
	}

	return res, nil
}

func (s *Search) Order(key string, lef int, rig int) ([]string, error) {
	con := s.pool.Get()
	defer con.Close()

	if lef < 0 {
		return nil, tracer.Maskf(executionFailedError, "lef must at least be 0")
	}

	if rig == 0 {
		return nil, tracer.Maskf(executionFailedError, "rig must not be 0")
	}
	if rig < -1 {
		return nil, tracer.Maskf(executionFailedError, "rig must at least be -1")
	}

	if rig != -1 && lef >= rig {
		return nil, tracer.Maskf(executionFailedError, "lef must be smaller than rig")
	}

	// Redis interprets the boundaries as inclusive numbers. We want to have
	// absolut numbers, because the second argument provided is about the
	// maximum number of elements. In case you want to have 1 element, providing
	// zero in this context would not make sense. Therefore we decrement all
	// numbers that are greater than zero. The exception is -1 which redis uses
	// to return all known elements. We want to keep this detail for our own
	// interface. So in case the user provides -1, we simply use it as is.
	if rig != -1 {
		rig--
	}

	res, err := redis.Strings(con.Do("ZREVRANGE", prefix.WithKeys(s.prefix, key), lef, rig))
	if err != nil {
		return nil, tracer.Mask(err)
	}

	return res, nil
}

func (s *Search) Score(key string, lef float64, rig float64) ([]string, error) {
	con := s.pool.Get()
	defer con.Close()

	res, err := redis.Strings(con.Do("ZREVRANGEBYSCORE", prefix.WithKeys(s.prefix, key), lef, rig))
	if err != nil {
		return nil, tracer.Mask(err)
	}

	return res, nil
}
