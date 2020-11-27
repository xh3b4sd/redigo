package sorted

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

type Create struct {
	pool *redis.Pool

	// elementScript is internally used to keep track of a script caching
	// mechanism when creating elements in sorted sets.
	elementScript *redis.Script

	prefix string
}

func (c *Create) Element(key string, val string, sco float64, ind ...string) error {
	con := c.pool.Get()
	defer con.Close()

	if c.elementScript == nil {
		scr := `
			if (ARGV[3] ~= nil) then
				# We got at least one index to keep track of. The first thing
				# we need to ensure is to verify any index we received does not
				# yet exist. As soon as we find a given index is already taken
				# we stop processing the request.
				local i = 3
				while ARGV[i] do
					local res = redis.call("SISMEMBER", KEYS[2], ARGV[i])
					if (res == 1) then
						return 0
					end

					i=i+1
				end

				# Only if we ensured that all indizes are not yet recorded we
				# can actually add them to our record. Tracking the indices
				# here aligns with the data persisted in the sorted set below.
				local j = 3
				while ARGV[j] do
					redis.call("SADD", KEYS[2], ARGV[j])

					j=j+1
				end
			end

			local val = ""
			local res = redis.call("ZRANGEBYSCORE", KEYS[1], ARGV[2], ARGV[2])
			for k, v in pairs(res) do
				val = v
				break
			end

			if (val ~= "") then
				return 1
			end

			redis.call("ZADD", KEYS[1], ARGV[2], ARGV[1])

			return 2
		`

		c.elementScript = redis.NewScript(2, scr)
	}

	var arg []interface{}
	{
		arg = append(arg, withPrefix(c.prefix, key))
		arg = append(arg, withPrefix(c.prefix, fmt.Sprintf("%s:ind", key)))
		arg = append(arg, val)
		arg = append(arg, sco)
		for _, s := range ind {
			arg = append(arg, s)
		}
	}

	res, err := redis.Int(c.elementScript.Do(con, arg...))
	if err != nil {
		return tracer.Mask(err)
	}

	switch res {
	case 0:
		return tracer.Maskf(alreadyExistsError, "index must be unique")
	case 1:
		return tracer.Maskf(alreadyExistsError, "score must be unique")
	case 2:
		return nil
	}

	return tracer.Mask(executionFailedError)
}
