package sorted

import (
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/redigo/index"
	"github.com/xh3b4sd/redigo/prefix"
	"github.com/xh3b4sd/tracer"
)

const searchIndexScript = `
	local sco = redis.call("ZSCORE", KEYS[2], ARGV[1])

	if (sco ~= false) then
		local res = redis.call("ZRANGE", KEYS[1], sco, sco, "BYSCORE")
		return res[1]
	end

	return ""
`

type search struct {
	pool *redis.Pool

	searchIndexScript *redis.Script

	prefix string
}

func (s *search) Index(key string, ind string) (string, error) {
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

	res, err := redis.String(s.searchIndexScript.Do(con, arg...))
	if err != nil {
		return "", tracer.Mask(err)
	}

	return res, nil
}

func (s *search) Inter(key ...string) ([]string, error) {
	con := s.pool.Get()
	defer con.Close()

	var mul []interface{}
	{
		mul = append(mul, len(key))

		for _, x := range key {
			mul = append(mul, prefix.WithKeys(s.prefix, x))
		}
	}

	res, err := redis.Strings(con.Do("ZINTER", mul...))
	if err != nil {
		return nil, tracer.Mask(err)
	}

	return res, nil
}

func (s *search) Order(key string, lef int, rig int, sco ...bool) ([]string, error) {
	con := s.pool.Get()
	defer con.Close()

	if len(sco) == 1 && !sco[0] {
		return nil, tracer.Maskf(executionFailedError, "sco must be true")
	}
	if len(sco) > 1 {
		return nil, tracer.Maskf(executionFailedError, "sco must not be provided more than once")
	}

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(s.prefix, key), lef, rig)

		if len(sco) == 1 {
			arg = append(arg, "WITHSCORES")
		}
	}

	res, err := redis.Strings(con.Do("ZRANGE", arg...))
	if err != nil {
		return nil, tracer.Mask(err)
	}

	return res, nil
}

func (s *search) Rando(key string, cou ...uint) ([]string, error) {
	con := s.pool.Get()
	defer con.Close()

	if len(cou) < 1 {
		cou = append(cou, 1)
	}
	if len(cou) > 1 {
		return nil, tracer.Maskf(executionFailedError, "cou must not be provided more than once")
	}

	res, err := redis.Strings(con.Do("ZRANDMEMBER", prefix.WithKeys(s.prefix, key), cou[0]))
	if IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, tracer.Mask(err)
	}

	return res, nil
}

func (s *search) Score(key string, lef float64, rig float64) ([]string, error) {
	con := s.pool.Get()
	defer con.Close()

	res, err := redis.Strings(con.Do("ZRANGE", prefix.WithKeys(s.prefix, key), lef, rig, "BYSCORE"))
	if err != nil {
		return nil, tracer.Mask(err)
	}

	return res, nil
}

func (s *search) Union(key ...string) ([]string, error) {
	con := s.pool.Get()
	defer con.Close()

	var mul []interface{}
	{
		mul = append(mul, len(key))

		for _, x := range key {
			mul = append(mul, prefix.WithKeys(s.prefix, x))
		}

		mul = append(mul, "AGGREGATE", "MIN")
	}

	res, err := redis.Strings(con.Do("ZUNION", mul...))
	if err != nil {
		return nil, tracer.Mask(err)
	}

	return res, nil
}
