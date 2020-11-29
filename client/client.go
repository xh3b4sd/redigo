package client

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/pool"
	"github.com/xh3b4sd/redigo/simple"
	"github.com/xh3b4sd/redigo/sorted"
)

type Config struct {
	Address string
	Pool    *redis.Pool
	Prefix  string
}

type Client struct {
	pool   *redis.Pool
	scored redigo.Sorted
	simple redigo.Simple
}

func New(config Config) (*Client, error) {
	if config.Address == "" {
		config.Address = "127.0.0.1:6379"
	}
	if config.Pool == nil {
		config.Pool = pool.NewPoolWithAddress(config.Address)
	}

	var err error

	var newSimple redigo.Simple
	{
		c := simple.Config{
			Pool: config.Pool,

			Prefix: config.Prefix,
		}

		newSimple, err = simple.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var newScored redigo.Sorted
	{
		c := sorted.Config{
			Pool: config.Pool,

			Prefix: config.Prefix,
		}

		newScored, err = sorted.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	c := &Client{
		pool:   config.Pool,
		scored: newScored,
		simple: newSimple,
	}

	return c, nil
}

func (c *Client) Check() error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("PING")
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (c *Client) Close() error {
	err := c.pool.Close()
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (c *Client) Purge() error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("FLUSHALL")
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (c *Client) Sorted() redigo.Sorted {
	return c.scored
}

func (c *Client) Simple() redigo.Simple {
	return c.simple
}
