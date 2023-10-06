package sorted

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo/pkg/index"
	"github.com/xh3b4sd/redigo/pkg/prefix"
)

const updateIndexScript = `
	-- Verify if the sorted set does even exist. We must not proceed if we would
	-- create a new sorted set and a new element within it.
	local exi = redis.call("EXISTS", KEYS[1])
	if (exi == 0) then
		return 0
	end

	local function upd(key, new, sco)
		-- We actually verified the existence of the element already. Now we
		-- only fetch the old value in order to perform the update.
		local res = redis.call("ZRANGE", key, sco, sco, "BYSCORE")
		local old = res[1]

		-- If the value did not change it might mean the indices did not change.
		-- We are ok with that internally. It is only important that the user
		-- facing element is properly reported to be updated or not.
		if (old == new) then
			return 2
		end

		redis.call("ZADD", key, sco, new)
		redis.call("ZREM", key, old)

		return 3
	end

	local function ver(key, new, sco)
		-- Verify if the score does already exist. If there is no element we
		-- cannot update it.
		local res = redis.call("ZRANGE", key, sco, sco, "BYSCORE")
		local old = res[1]
		if (old == nil) then
			return 1
		end

		-- Verify if the existing value is already what we want to update to. If
		-- the desired state is already reconciled we do not need to proceed
		-- further.
		if (old == new) then
			return 2
		end

		return 3
	end

	-- Verify all scores have associated values. We need to do this upfront for
	-- the given element and the internally managed indices.
	local i = 3
	while ARGV[i] do
		local res = ver(KEYS[2], ARGV[i], ARGV[2])
		if (res == 1) then
			return res
		end

		i=i+1
	end
	local res = ver(KEYS[1], ARGV[1], ARGV[2])
	if (res == 1) then
		return res
	end

	-- Only if all verifications are completed successfully and there is no
	-- reason to fail anymore we can continue to actually process the updates.
	local j = 3
	while ARGV[j] do
		upd(KEYS[2], ARGV[j], ARGV[2])

		j=j+1
	end

	return upd(KEYS[1], ARGV[1], ARGV[2])
`

const updateScoreScript = `
	-- Verify if the score does already exist. If there is no element we
	-- cannot update it.
	local res = redis.call("ZRANGE", KEYS[1], ARGV[2], ARGV[2], "BYSCORE")

	local old = res[1]
	if (old == nil) then
		return 0
	end

	-- Verify if the existing value is already what we want to update to. If
	-- the desired state is already reconciled we do not need to proceed
	-- further.
	if (old == ARGV[1]) then
		return 1
	end

	redis.call("ZADD", KEYS[1], ARGV[2], ARGV[1])
	redis.call("ZREM", KEYS[1], old)

	return 2
`

type update struct {
	pool *redis.Pool

	updateIndexScript *redis.Script
	updateScoreScript *redis.Script

	prefix string
}

// Index executes a script of three key operations in order to reliably modify
// the element of a sorted set. Consider the element bef being created using the
// score 23 like shown below.
//
//	redis> ZADD k:foo 25 "old"
//	(integer) 1
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
//	redis> ZRANGE k:foo 25 25 BYSCORE
//	1) "old"
//
//	redis> ZADD k:foo 25 "new"
//	(integer) 1
//
//	redis> ZREM k:foo "old"
//	(integer) 1
func (u *update) Index(key string, new string, sco float64, ind ...string) (bool, error) {
	con := u.pool.Get()
	defer con.Close()

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(u.prefix, key))            // KEYS[1]
		arg = append(arg, prefix.WithKeys(u.prefix, index.New(key))) // KEYS[2]
		arg = append(arg, new)                                       // ARGV[1]
		arg = append(arg, sco)                                       // ARGV[2]

		for _, s := range ind {
			arg = append(arg, s)
		}
	}

	res, err := redis.Int(u.updateIndexScript.Do(con, arg...))
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

func (u *update) Score(key string, new string, sco float64) (bool, error) {
	con := u.pool.Get()
	defer con.Close()

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(u.prefix, key)) // KEYS[1]
		arg = append(arg, new)                            // ARGV[1]
		arg = append(arg, sco)                            // ARGV[2]
	}

	res, err := redis.Int(u.updateScoreScript.Do(con, arg...))
	if err != nil {
		return false, tracer.Mask(err)
	}

	switch res {
	case 0:
		return false, tracer.Maskf(notFoundError, "element does not exist in sorted set")
	case 1:
		return false, nil
	case 2:
		return true, nil
	}

	return false, tracer.Mask(executionFailedError)
}
