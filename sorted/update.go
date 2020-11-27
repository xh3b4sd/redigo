package sorted

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

type Update struct {
	pool *redis.Pool

	// updateScript is internally used to keep track of a script caching
	// mechanism when updating elements in sorted sets.
	updateScript *redis.Script

	prefix string
}

// Value executes a script of three key operations in order to reliably modify
// the element of a sorted set. Consider the element bef being created using the
// score 23 like shown below.
//
//     redis> ZADD k:foo 25 "old"
//     (integer) 1
//
// The first step being executed is to lookup the old element associated with
// score. This is necessary because removing an element from a sorted set in our
// case does work best knowing the element itself. That way we can add the new
// element first assuming that in worst case we have two elements instead of
// zero.
//
// The second step being executed is to add the new element associated with
// score. At this point we have two elements. The old and the new one.
//
// The third step being executed is to remove the old element associated with
// score. Now the value retrieved in the first step is being leveraged. Removing
// the old element marks the end of the executed transaction, leaving a clean
// state of an updated element behind.
//
//     redis> ZREVRANGE k:foo 25 25
//     1) "old"
//
//     redis> ZADD k:foo 25 "new"
//     (integer) 1
//
//     redis> ZREM k:foo "old"
//     (integer) 1
//
func (u *Update) Value(key string, new string, sco float64) (bool, error) {
	con := u.pool.Get()
	defer con.Close()

	if u.updateScript == nil {
		scr := `
			local exi = redis.call("EXISTS", KEYS[1])
			if (exi == 0) then
				return 0
			end

			local old = ""
			local res = redis.call("ZRANGEBYSCORE", KEYS[1], ARGV[2], ARGV[2])
			for k, v in pairs(res) do
				old = v
				break
			end

			if (old == "") then
				return 1
			end

			if (old == ARGV[1]) then
				return 2
			end

			redis.call("ZADD", KEYS[1], ARGV[2], ARGV[1])
			redis.call("ZREM", KEYS[1], old)

			return 3
		`

		u.updateScript = redis.NewScript(1, scr)
	}

	res, err := redis.Int(u.updateScript.Do(con, withPrefix(u.prefix, key), new, sco))
	if err != nil {
		return false, tracer.Mask(err)
	}

	switch res {
	case 0:
		return false, tracer.Maskf(notFoundError, "sorted set does not exist under key")
	case 1:
		return false, tracer.Maskf(notFoundError, "element does not exist in sorted set")
	case 2:
		return false, nil
	case 3:
		return true, nil
	}

	return false, tracer.Mask(executionFailedError)
}
