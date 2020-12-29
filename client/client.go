package client

import (
	"net"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/pool"
	"github.com/xh3b4sd/redigo/simple"
	"github.com/xh3b4sd/redigo/sorted"
)

const (
	KindSimple   = "simple"
	KindSentinel = "sentinel"
)

type Config struct {
	Address string
	Kind    string
	Pool    *redis.Pool
	Prefix  string
}

type Client struct {
	pool   *redis.Pool
	scored redigo.Sorted
	simple redigo.Simple
}

func New(config Config) (*Client, error) {
	if config.Kind != KindSimple && config.Kind != KindSentinel {
		return nil, tracer.Maskf(invalidConfigError, "%T.Kind must be %s or %s", config, KindSimple, KindSentinel)
	}

	if config.Address == "" && config.Kind == KindSimple {
		config.Address = defaultSimpleAddress()
	}
	if config.Address == "" && config.Kind == KindSentinel {
		config.Address = defaultSentinelAddress()
	}

	if config.Pool == nil && config.Kind == KindSimple {
		config.Pool = pool.NewSimplePoolWithAddress(config.Address)
	}
	if config.Pool == nil && config.Kind == KindSentinel {
		config.Pool = pool.NewSentinelPoolWithAddress(config.Address)
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

func defaultSimpleAddress() string {
	var hos string
	{
		hos = os.Getenv("REDIS_HOST")
		if hos == "" {
			hos = "127.0.0.1"
		}
	}

	var por string
	{
		por = os.Getenv("REDIS_PORT")
		if por == "" {
			por = "6379"
		}
	}

	return net.JoinHostPort(hos, por)
}

func defaultSentinelAddress() string {
	var hos string
	{
		hos = os.Getenv("REDIS_SENTINEL_HOST")
		if hos == "" {
			hos = "127.0.0.1"
		}
	}

	var por string
	{
		por = os.Getenv("REDIS_SENTINEL_PORT")
		if por == "" {
			por = "26379"
		}
	}

	return net.JoinHostPort(hos, por)
}
