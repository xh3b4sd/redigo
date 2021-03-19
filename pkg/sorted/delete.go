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

const deleteScoreScript = `
	redis.call("ZREMRANGEBYSCORE", KEYS[2], ARGV[1], ARGV[1])
	redis.call("ZREMRANGEBYSCORE", KEYS[1], ARGV[1], ARGV[1])

	return 0
`

const deleteValueScript = `
	local sco = redis.call("ZSCORE", KEYS[1], ARGV[1])

	if (sco ~= false) then
		redis.call("ZREMRANGEBYSCORE", KEYS[2], sco, sco)
		redis.call("ZREMRANGEBYSCORE", KEYS[1], sco, sco)
	else
		redis.call("ZREM", KEYS[1], ARGV[1])
	end

	return 0
`

type Delete struct {
	pool *redis.Pool

	deleteCleanScript *redis.Script
	deleteScoreScript *redis.Script
	deleteValueScript *redis.Script

	prefix string
}

func (d *Delete) Clean(key string) error {
	con := d.pool.Get()
	defer con.Close()

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(d.prefix, key))
		arg = append(arg, prefix.WithKeys(d.prefix, index.New(key)))
	}

	_, err := redis.Int(d.deleteCleanScript.Do(con, arg...))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (d *Delete) Score(key string, sco float64) error {
	con := d.pool.Get()
	defer con.Close()

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(d.prefix, key))
		arg = append(arg, prefix.WithKeys(d.prefix, index.New(key)))
		arg = append(arg, sco)
	}

	_, err := redis.Int(d.deleteScoreScript.Do(con, arg...))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (d *Delete) Value(key string, val string) error {
	con := d.pool.Get()
	defer con.Close()

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(d.prefix, key))
		arg = append(arg, prefix.WithKeys(d.prefix, index.New(key)))
		arg = append(arg, val)
	}

	_, err := redis.Int(d.deleteValueScript.Do(con, arg...))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
