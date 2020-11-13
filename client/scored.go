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

func (s *Scored) CutOff(key string, num int) error {
	conn := s.pool.Get()
	defer conn.Close()

	length, err := redis.Int(conn.Do("ZCARD", withPrefix(s.prefix, key)))
	if err != nil {
		return tracer.Mask(err)
	}

	count := length - num
	if count < 1 {
		return nil
	}

	_, err = redis.Strings(conn.Do("ZPOPMIN", withPrefix(s.prefix, key), count))
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

	if lef >= rig {
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
