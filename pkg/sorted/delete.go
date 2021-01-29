package sorted

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo/pkg/prefix"
)

const deleteScoreScript = `
	-- Ensure that all the indizes we have recorded get deleted. Deleting by
	-- score is easy because we can delete all indizes by score at once.
	redis.call("ZREMRANGEBYSCORE", KEYS[2], ARGV[1], ARGV[1])

	redis.call("ZREMRANGEBYSCORE", KEYS[1], ARGV[1], ARGV[1])

	return 0
`

const deleteValueScript = `
	if (ARGV[2] ~= nil) then
		-- Ensure that all the indizes we have recorded get deleted. Deleting by
		-- value is complex because we have to delete all indizes by value each.
		local j = 2
		while ARGV[j] do
			redis.call("ZREM", KEYS[2], ARGV[j])

			j=j+1
		end
	end

	redis.call("ZREM", KEYS[1], ARGV[1])

	return 0
`

type Delete struct {
	pool *redis.Pool

	scoreScript *redis.Script
	valueScript *redis.Script

	prefix string
}

func (d *Delete) Score(key string, sco float64) error {
	con := d.pool.Get()
	defer con.Close()

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(d.prefix, key))
		arg = append(arg, prefix.WithKeys(d.prefix, fmt.Sprintf("%s:ind", key)))
		arg = append(arg, sco)
	}

	_, err := redis.Int(d.scoreScript.Do(con, arg...))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (d *Delete) Value(key string, val string, ind ...string) error {
	con := d.pool.Get()
	defer con.Close()

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(d.prefix, key))
		arg = append(arg, prefix.WithKeys(d.prefix, fmt.Sprintf("%s:ind", key)))
		arg = append(arg, val)
		for _, s := range ind {
			arg = append(arg, s)
		}
	}

	_, err := redis.Int(d.valueScript.Do(con, arg...))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
