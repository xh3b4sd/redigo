package client

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

type Scored struct {
	pool *redis.Pool

	prefix string
}

func (s *Scored) Create(key string, element string, score float64) error {
	conn := s.pool.Get()
	defer conn.Close()

	_, err := redis.Int(conn.Do("ZADD", withPrefix(s.prefix, key), score, element))
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

func (s *Scored) Delete(key string, element string) error {
	conn := s.pool.Get()
	defer conn.Close()

	_, err := redis.Int(conn.Do("ZREM", withPrefix(s.prefix, key), element))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

// Search returns the list of scored elements stored under key. Note that num
// may be -1 in order to list all elements.
func (s *Scored) Search(key string, num int) ([]string, error) {
	conn := s.pool.Get()
	defer conn.Close()

	if num == 0 {
		return nil, tracer.Maskf(executionFailedError, "num must not be 0")
	}
	if num < -1 {
		return nil, tracer.Maskf(executionFailedError, "num must at least be -1")
	}

	// Redis interprets the boundaries as inclusive numbers. We want to have
	// absolut numbers, because the second argument provided is about the maximum
	// number of elements. In case you want to have 1 element, providing zero in
	// this context would not make sense. Therefore we decrement all numbers that
	// are greater than zero.
	num--

	result, err := redis.Strings(conn.Do("ZREVRANGE", withPrefix(s.prefix, key), 0, num, "WITHSCORES"))
	if err != nil {
		return nil, tracer.Mask(err)
	}

	return result, nil
}
