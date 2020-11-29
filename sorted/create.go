package sorted

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo/prefix"
)

const createElementScript = `
	-- Verify if the score does already exist. With this implementation the
	-- score is treated as ID and must therefore be unique.
	local val = ""
	local res = redis.call("ZRANGEBYSCORE", KEYS[1], ARGV[2], ARGV[2])
	for k, v in pairs(res) do
		val = v
		break
	end
	if (val ~= "") then
		return 0
	end

	if (ARGV[3] ~= nil) then
		-- We got at least one index to keep track of. The first thing we need
		-- to ensure is to verify any index we received does not yet exist. As
		-- soon as we find a given index is already taken we stop processing the
		-- request.
		local i = 3
		while ARGV[i] do
			local res = redis.call("ZSCORE", KEYS[2], ARGV[i])
			if (res ~= nil) then
				return 1
			end

			i=i+1
		end

		-- Only if we ensured that all indizes are not yet recorded we can
		-- actually add them to our record. Tracking the indices here aligns
		-- with the data persisted in the sorted set below.
		local j = 3
		while ARGV[j] do
			redis.call("ZADD", KEYS[2], ARGV[2], ARGV[j])

			j=j+1
		end
	end

	redis.call("ZADD", KEYS[1], ARGV[2], ARGV[1])

	return 2
`

type Create struct {
	pool *redis.Pool

	createElementScript *redis.Script

	prefix string
}

func (c *Create) Element(key string, val string, sco float64, ind ...string) error {
	con := c.pool.Get()
	defer con.Close()

	if c.createElementScript == nil {
		c.createElementScript = redis.NewScript(2, createElementScript)
	}

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(c.prefix, key))
		arg = append(arg, prefix.WithKeys(c.prefix, fmt.Sprintf("%s:ind", key)))
		arg = append(arg, val)
		arg = append(arg, sco)
		for _, s := range ind {
			arg = append(arg, s)
		}
	}

	res, err := redis.Int(c.createElementScript.Do(con, arg...))
	if err != nil {
		return tracer.Mask(err)
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
