package client

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

// PoolConfig represents the configuration used to create a new redis
// pool.
type PoolConfig struct {
	// Dial is an application supplied function for creating and configuring a
	// redis connection on demand.
	Dial func() (redis.Conn, error)
	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration
	// MaxIdle is the allowed maximum number of idle connections in the pool.
	MaxIdle int
	// MaxActive is the allowed maximum number of connections allocated by the
	// pool at a given time.  When zero, there is no limit on the number of
	// connections in the pool.
	MaxActive int
	// TestOnBorrow is an optional application supplied function for checking the
	// health of an idle connection before the connection is used again by the
	// application. Argument t is the time that the connection was returned to the
	// pool. If the function returns an error, then the connection is closed.
	TestOnBorrow func(c redis.Conn, t time.Time) error
}

// DefaultPoolConfig provides a default configuration to create a new
// redis pool by best effort.
func DefaultPoolConfig() PoolConfig {
	newConfig := PoolConfig{
		MaxIdle:     100,
		MaxActive:   100,
		IdleTimeout: 5 * time.Minute,
		Dial:        NewDial(DefaultDialConfig()),
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

	return newConfig
}

// NewPool creates a new configured redis pool.
func NewPool(config PoolConfig) *redis.Pool {
	p := &redis.Pool{
		Dial:         config.Dial,
		IdleTimeout:  config.IdleTimeout,
		MaxIdle:      config.MaxIdle,
		MaxActive:    config.MaxActive,
		TestOnBorrow: config.TestOnBorrow,
	}

	return p
}

func NewPoolWithAddress(address string) *redis.Pool {
	newDialConfig := DefaultDialConfig()
	newDialConfig.Address = address

	var p *redis.Pool
	{
		c := DefaultPoolConfig()
		c.Dial = NewDial(newDialConfig)

		p = NewPool(c)
	}

	return p
}
