package pool

import (
	"errors"
	"time"

	"github.com/FZambia/sentinel"
	"github.com/gomodule/redigo/redis"
)

func NewSimplePoolWithAddress(address string) *redis.Pool {
	var p *redis.Pool
	{
		p = &redis.Pool{
			MaxIdle:     100,
			MaxActive:   100,
			IdleTimeout: 5 * time.Minute,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", address)
				if err != nil {
					return nil, err
				}

				return c, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}

				_, err := c.Do("PING")
				if err != nil {
					return err
				}

				return nil
			},
		}
	}

	return p
}

func NewSimplePoolWithConnection(connection redis.Conn) *redis.Pool {
	var p *redis.Pool
	{
		p = &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return connection, nil
			},
		}
	}

	return p
}

func NewSentinelPoolWithAddress(address string) *redis.Pool {
	timeout := 500 * time.Millisecond

	sntnl := &sentinel.Sentinel{
		Addrs:      []string{address},
		MasterName: "mymaster",
		Dial: func(addr string) (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				addr,
				redis.DialConnectTimeout(timeout),
				redis.DialReadTimeout(timeout),
				redis.DialWriteTimeout(timeout),
			)
			if err != nil {
				return nil, err
			}

			return c, nil
		},
	}

	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   64,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			a, err := sntnl.MasterAddr()
			if err != nil {
				return nil, err
			}

			c, err := redis.Dial("tcp", a)
			if err != nil {
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if !sentinel.TestRole(c, "master") {
				return errors.New("role check failed")
			} else {
				return nil
			}
		},
	}
}
