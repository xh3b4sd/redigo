package delete

import (
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
	cln *redis.Script
	ind *redis.Script
}

func New(c Config) *Redis {
	return &Redis{
		poo: c.Poo,
		pre: c.Pre,
		cln: redis.NewScript(2, deleteCleanScript),
		ind: redis.NewScript(2, deleteIndexScript),
	}
}

func (r *Redis) Clean(key string) error {
	var err error

	var con redis.Conn
	{
		con = r.poo.Get()
		defer con.Close()
	}

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(r.pre, key))            // KEYS[1]
		arg = append(arg, prefix.WithKeys(r.pre, index.New(key))) // KEYS[2]
	}

	{
		_, err = redis.Int64(r.cln.Do(con, arg...))
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

func (r *Redis) Index(key string, val ...string) error {
	var err error

	var con redis.Conn
	{
		con = r.poo.Get()
		defer con.Close()
	}

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(r.pre, key))            // KEYS[1]
		arg = append(arg, prefix.WithKeys(r.pre, index.New(key))) // KEYS[2]

		for _, x := range val {
			arg = append(arg, x)
		}
	}

	{
		_, err = redis.Int64(r.ind.Do(con, arg...))
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

func (r *Redis) Limit(key string, lim int) error {
	var err error

	var con redis.Conn
	{
		con = r.poo.Get()
		defer con.Close()
	}

	if lim < 0 {
		return tracer.Maskf(executionFailedError, "lim must at least be 0")
	}

	var sta int
	{
		sta = 0
	}

	var end int
	{
		end = -1 * (lim + 1)
	}

	{
		_, err = redis.Int64(con.Do("ZREMRANGEBYRANK", prefix.WithKeys(r.pre, key), sta, end))
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

func (r *Redis) Score(key string, sco float64, end ...float64) error {
	var err error

	var con redis.Conn
	{
		con = r.poo.Get()
		defer con.Close()
	}

	if len(end) > 1 {
		return tracer.Maskf(executionFailedError, "end must not be provided more than once")
	}

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(r.pre, key))

		if len(end) == 1 {
			arg = append(arg, sco, end[0])
		} else {
			arg = append(arg, sco, sco)
		}
	}

	{
		_, err = con.Do("ZREMRANGEBYSCORE", arg...)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

func (r *Redis) Value(key string, val ...string) error {
	var err error

	var con redis.Conn
	{
		con = r.poo.Get()
		defer con.Close()
	}

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(r.pre, key))

		for _, x := range val {
			arg = append(arg, x)
		}
	}

	{
		_, err = redis.Int64(con.Do("ZREM", arg...))
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
