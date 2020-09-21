package client

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

type Simple struct {
	pool *redis.Pool

	prefix string
}

func (s *Simple) Create(key, element string) error {
	conn := s.pool.Get()
	defer conn.Close()

	reply, err := redis.String(conn.Do("SET", withPrefix(s.prefix, key), element))
	if err != nil {
		return tracer.Mask(err)
	}

	if reply != "OK" {
		return tracer.Maskf(executionFailedError, "SET not executed correctly")
	}

	return nil
}

func (s *Simple) Delete(key string) error {
	conn := s.pool.Get()
	defer conn.Close()

	_, err := redis.Int64(conn.Do("DEL", withPrefix(s.prefix, key)))
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (s *Simple) Exists(key string) (bool, error) {
	conn := s.pool.Get()
	defer conn.Close()

	result, err := redis.Bool(conn.Do("EXISTS", withPrefix(s.prefix, key)))
	if err != nil {
		return false, tracer.Mask(err)
	}

	return result, nil
}

func (s *Simple) Search(key string) (string, error) {
	conn := s.pool.Get()
	defer conn.Close()

	result, err := redis.String(conn.Do("GET", withPrefix(s.prefix, key)))
	if IsNotFound(err) {
		return "", tracer.Maskf(notFoundError, key)
	} else if err != nil {
		return "", tracer.Mask(err)
	}

	return result, nil
}
