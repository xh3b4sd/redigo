package client

import (
	"sync"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo"
)

type Config struct {
	Address string
	Pool    *redis.Pool
	Prefix  string
}

type Client struct {
	pool         *redis.Pool
	shutdownOnce sync.Once
	scored       redigo.Scored
	simple       redigo.Simple
}

func New(config Config) (*Client, error) {
	if config.Address == "" {
		config.Address = "127.0.0.1:6379"
	}
	if config.Pool == nil {
		config.Pool = NewPoolWithAddress(config.Address)
	}

	var newSimple redigo.Simple
	{
		newSimple = &Simple{
			pool: config.Pool,

			prefix: config.Prefix,
		}
	}

	var newScored redigo.Scored
	{
		newScored = &Scored{
			pool: config.Pool,

			prefix: config.Prefix,
		}
	}

	c := &Client{
		pool:         config.Pool,
		shutdownOnce: sync.Once{},
		scored:       newScored,
		simple:       newSimple,
	}

	return c, nil
}

func (c *Client) Ping() error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("PING")
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (c *Client) Scored() redigo.Scored {
	return c.scored
}

func (c *Client) Shutdown() {
	c.shutdownOnce.Do(func() {
		c.pool.Close()
	})
}

func (c *Client) Simple() redigo.Simple {
	return c.simple
}

func withPrefix(prefix string, keys ...string) string {
	newKey := prefix

	for _, k := range keys {
		newKey += ":" + k
	}

	if prefix == "" {
		newKey = newKey[1:]
	}

	return newKey
}
