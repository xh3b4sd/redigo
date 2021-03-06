package sorted

import (
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo/pkg/index"
	"github.com/xh3b4sd/redigo/pkg/prefix"
)

const createElementScript = `
	-- Verify if the score does already exist. The first key here might be "ssk"
	-- and the second argument might be "0.8". If we get any value in response
	-- the score is already taken.
	local res = redis.call("ZRANGEBYSCORE", KEYS[1], ARGV[2], ARGV[2])
	if (res[1] ~= nil) then
		return 0
	end

	if (ARGV[3] ~= nil) then
		-- Verify if the index does already exist. The second key here might be
		-- "ssk:ind" and the argument might be "name". If we get any value in
		-- response the index is already taken.
		local i = 3
		while ARGV[i] do
			local res = redis.call("ZSCORE", KEYS[2], ARGV[i])
			if (res ~= false) then
				return 1
			end

			i=i+1
		end

		-- Only if we ensured that the score is unique and that all indizes are
		-- not yet recorded, we can then add them to our sorted sets.
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

	if len(ind) != 0 {
		m := map[string]int{}
		for _, s := range ind {
			m[s] = m[s] + 1
		}

		for _, v := range m {
			if v > 1 {
				return tracer.Maskf(executionFailedError, "index must be unique")
			}
		}

		for _, s := range ind {
			if s == "" {
				return tracer.Maskf(executionFailedError, "index must not be empty")
			}
			if strings.Count(s, " ") != 0 {
				return tracer.Maskf(executionFailedError, "index must not contain whitespace")
			}
		}
	}

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(c.prefix, key))
		arg = append(arg, prefix.WithKeys(c.prefix, index.New(key)))
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
