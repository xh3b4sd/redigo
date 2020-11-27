package client

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

type Scored struct {
	pool *redis.Pool

	// createScript is internally used to keep track of a script caching
	// mechanism when creating elements in sorted sets.
	createScript *redis.Script
	// updateScript is internally used to keep track of a script caching
	// mechanism when updating elements in sorted sets.
	updateScript *redis.Script

	prefix string
}

func (s *Scored) Create(key string, ele string, sco float64) error {
	con := s.pool.Get()
	defer con.Close()

	if s.createScript == nil {
		scr := `
			local old = ""
			local res = redis.call("ZRANGEBYSCORE", KEYS[1], ARGV[2], ARGV[2])
			for k, v in pairs(res) do
				old = v
				break
			end

			if (old != "") then
				return 0
			end

			redis.call("ZADD", KEYS[1], ARGV[2], ARGV[1])

			return 1
		`

		s.createScript = redis.NewScript(1, scr)
	}

	res, err := redis.Int(s.createScript.Do(con, withPrefix(s.prefix, key), ele, sco))
	if err != nil {
		return tracer.Mask(err)
	}

	switch res {
	case 0:
		return tracer.Maskf(alreadyExistsError, "score must be unique")
	case 1:
		return nil
	}

	return tracer.Mask(executionFailedError)
}

func (s *Scored) Delete(key string, ele string) error {
	con := s.pool.Get()
	defer con.Close()

	_, err := redis.Int(con.Do("ZREM", withPrefix(s.prefix, key), ele))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (s *Scored) Exists(key string) (bool, error) {
	con := s.pool.Get()
	defer con.Close()

	result, err := redis.Bool(con.Do("EXISTS", withPrefix(s.prefix, key)))
	if err != nil {
		return false, tracer.Mask(err)
	}

	return result, nil
}

func (s *Scored) Search(key string, lef int, rig int) ([]string, error) {
	con := s.pool.Get()
	defer con.Close()

	if lef < 0 {
		return nil, tracer.Maskf(executionFailedError, "lef must at least be 0")
	}

	if rig == 0 {
		return nil, tracer.Maskf(executionFailedError, "rig must not be 0")
	}
	if rig < -1 {
		return nil, tracer.Maskf(executionFailedError, "rig must at least be -1")
	}

	if rig != -1 && lef >= rig {
		return nil, tracer.Maskf(executionFailedError, "lef must be smaller than rig")
	}

	// Redis interprets the boundaries as inclusive numbers. We want to have
	// absolut numbers, because the second argument provided is about the
	// maximum number of elements. In case you want to have 1 element, providing
	// zero in this context would not make sense. Therefore we decrement all
	// numbers that are greater than zero. The exception is -1 which redis uses
	// to return all known elements. We want to keep this detail for our own
	// interface. So in case the user provides -1, we simply use it as is.
	if rig != -1 {
		rig--
	}

	result, err := redis.Strings(con.Do("ZREVRANGE", withPrefix(s.prefix, key), lef, rig))
	if err != nil {
		return nil, tracer.Mask(err)
	}

	return result, nil
}

// Update executes a script of three key operations in order to reliably modify
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
func (s *Scored) Update(key string, new string, sco float64) (bool, error) {
	con := s.pool.Get()
	defer con.Close()

	if s.updateScript == nil {
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

		s.updateScript = redis.NewScript(1, scr)
	}

	res, err := redis.Int(s.updateScript.Do(con, withPrefix(s.prefix, key), new, sco))
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
