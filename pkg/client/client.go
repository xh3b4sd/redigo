package client

import (
	"net"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/pkg/locker"
	"github.com/xh3b4sd/redigo/pkg/pool"
	"github.com/xh3b4sd/redigo/pkg/pubsub"
	"github.com/xh3b4sd/redigo/pkg/simple"
	"github.com/xh3b4sd/redigo/pkg/sorted"
	"github.com/xh3b4sd/redigo/pkg/walker"
)

const (
	KindSingle   = "single"
	KindSentinel = "sentinel"
)

type Config struct {
	Address string
	Count   int
	Kind    string
	Pool    *redis.Pool
	Prefix  string
}

type Client struct {
	pool *redis.Pool

	locker redigo.Locker
	pubSub redigo.PubSub
	simple redigo.Simple
	sorted redigo.Sorted
	walker redigo.Walker
}

func New(config Config) (*Client, error) {
	if config.Kind != KindSingle && config.Kind != KindSentinel {
		return nil, tracer.Maskf(invalidConfigError, "%T.Kind must be %s or %s", config, KindSingle, KindSentinel)
	}

	if config.Address == "" && config.Kind == KindSingle {
		config.Address = defaultSingleAddress()
	}
	if config.Address == "" && config.Kind == KindSentinel {
		config.Address = defaultSentinelAddress()
	}

	if config.Pool == nil && config.Kind == KindSingle {
		config.Pool = pool.NewSinglePoolWithAddress(config.Address)
	}
	if config.Pool == nil && config.Kind == KindSentinel {
		config.Pool = pool.NewSentinelPoolWithAddress(config.Address)
	}

	var err error

	var newLocker redigo.Locker
	{
		c := locker.Config{
			Pool: config.Pool,

			Prefix: config.Prefix,
		}

		newLocker, err = locker.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var newPubSub redigo.PubSub
	{
		c := pubsub.Config{
			Pool: config.Pool,

			Prefix: config.Prefix,
		}

		newPubSub, err = pubsub.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

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

	var newSorted redigo.Sorted
	{
		c := sorted.Config{
			Pool: config.Pool,

			Prefix: config.Prefix,
		}

		newSorted, err = sorted.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var newWalker redigo.Walker
	{
		c := walker.Config{
			Pool: config.Pool,

			Count:  config.Count,
			Prefix: config.Prefix,
		}

		newWalker, err = walker.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	c := &Client{
		pool: config.Pool,

		locker: newLocker,
		pubSub: newPubSub,
		simple: newSimple,
		sorted: newSorted,
		walker: newWalker,
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

func (c *Client) Locker() redigo.Locker {
	return c.locker
}

func (c *Client) PubSub() redigo.PubSub {
	return c.pubSub
}

func (c *Client) Simple() redigo.Simple {
	return c.simple
}

func (c *Client) Sorted() redigo.Sorted {
	return c.sorted
}

func (c *Client) Walker() redigo.Walker {
	return c.walker
}

func defaultSingleAddress() string {
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
