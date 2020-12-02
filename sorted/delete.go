package sorted

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo/prefix"
)

const deleteElementScript = `
	if (ARGV[2] ~= nil) then
		-- Only if we ensured that all indizes are not yet recorded we can
		-- actually add them to our record. Tracking the indices here aligns
		-- with the data persisted in the sorted set below.
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

	elementScript *redis.Script

	prefix string
}

func (d *Delete) Element(key string, val string, ind ...string) error {
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

	_, err := redis.Int(d.elementScript.Do(con, arg...))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
