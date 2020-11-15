package client

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

type Scored struct {
	pool *redis.Pool

	prefix string
}

func (s *Scored) Create(key string, ele string, sco float64) error {
	conn := s.pool.Get()
	defer conn.Close()

	_, err := redis.Int(conn.Do("ZADD", withPrefix(s.prefix, key), sco, ele))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (s *Scored) Delete(key string, ele string) error {
	conn := s.pool.Get()
	defer conn.Close()

	_, err := redis.Int(conn.Do("ZREM", withPrefix(s.prefix, key), ele))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (s *Scored) Search(key string, lef int, rig int) ([]string, error) {
	conn := s.pool.Get()
	defer conn.Close()

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

	result, err := redis.Strings(conn.Do("ZREVRANGE", withPrefix(s.prefix, key), lef, rig))
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
func (s *Scored) Update(key string, new string, sco float64) error {
	conn := s.pool.Get()
	defer conn.Close()

	scr := `
        local old = redis.call("ZREVRANGE", KEYS[1], ARGV[2], ARGV[2])
        redis.call("ZADD", KEYS[1], ARGV[2], ARGV[1])
        redis.call("ZREM", KEYS[1], old)
	`

	_, err := conn.Do("EVAL", scr, 1, withPrefix(s.prefix, key), new, sco)
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
