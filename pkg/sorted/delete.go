package sorted

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo/pkg/index"
	"github.com/xh3b4sd/redigo/pkg/prefix"
)

const deleteCleanScript = `
	redis.call("DEL", KEYS[2])
	redis.call("DEL", KEYS[1])

	return 0
`

const deleteIndexScript = `
	for i = 1, #ARGV do
			local sco = redis.call("ZSCORE", KEYS[1], ARGV[i])

			if (sco ~= false) then
					redis.call("ZREMRANGEBYSCORE", KEYS[2], sco, sco)
			end
	end

	redis.call("ZREM", KEYS[1], unpack(ARGV))

	return 0
`

const deleteScoreScript = `
	redis.call("ZREMRANGEBYSCORE", KEYS[2], ARGV[1], ARGV[1])
	redis.call("ZREMRANGEBYSCORE", KEYS[1], ARGV[1], ARGV[1])

	return 0
`

type delete struct {
	pool *redis.Pool

	deleteCleanScript *redis.Script
	deleteIndexScript *redis.Script
	deleteScoreScript *redis.Script

	prefix string
}

func (d *delete) Clean(key string) error {
	con := d.pool.Get()
	defer con.Close()

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(d.prefix, key))            // KEYS[1]
		arg = append(arg, prefix.WithKeys(d.prefix, index.New(key))) // KEYS[2]
	}

	_, err := redis.Int(d.deleteCleanScript.Do(con, arg...))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (d *delete) Index(key string, val ...string) error {
	con := d.pool.Get()
	defer con.Close()

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(d.prefix, key))            // KEYS[1]
		arg = append(arg, prefix.WithKeys(d.prefix, index.New(key))) // KEYS[2]

		for _, x := range val {
			arg = append(arg, x)
		}
	}

	_, err := redis.Int(d.deleteIndexScript.Do(con, arg...))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (d *delete) Limit(key string, lim int) error {
	con := d.pool.Get()
	defer con.Close()

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

	_, err := redis.Int64(con.Do("ZREMRANGEBYRANK", prefix.WithKeys(d.prefix, key), sta, end))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (d *delete) Score(key string, sco float64) error {
	con := d.pool.Get()
	defer con.Close()

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(d.prefix, key))            // KEYS[1]
		arg = append(arg, prefix.WithKeys(d.prefix, index.New(key))) // KEYS[2]
		arg = append(arg, sco)                                       // ARGV[1]
	}

	_, err := redis.Int(d.deleteScoreScript.Do(con, arg...))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (d *delete) Value(key string, val ...string) error {
	con := d.pool.Get()
	defer con.Close()

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(d.prefix, key))

		for _, x := range val {
			arg = append(arg, x)
		}
	}

	_, err := redis.Int64(con.Do("ZREM", arg...))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
