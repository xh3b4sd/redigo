package client

import (
	"net"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/budget/v3"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo/pkg/backup"
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
	Locker  ConfigLocker
	Pool    *redis.Pool
	Prefix  string
}

type ConfigLocker struct {
	Budget budget.Interface
	Expiry time.Duration
	Name   string
}

type Client struct {
	bac backup.Interface
	loc locker.Interface
	poo *redis.Pool
	pub pubsub.Interface
	sim simple.Interface
	sor sorted.Interface
	wal walker.Interface
}

func New(con Config) (*Client, error) {
	if con.Kind != KindSingle && con.Kind != KindSentinel {
		return nil, tracer.Maskf(invalidConfigError, "%T.Kind must be %s or %s", con, KindSingle, KindSentinel)
	}

	if con.Address == "" && con.Kind == KindSingle {
		con.Address = defaultSingleAddress()
	}
	if con.Address == "" && con.Kind == KindSentinel {
		con.Address = defaultSentinelAddress()
	}

	if con.Pool == nil && con.Kind == KindSingle {
		con.Pool = pool.NewSinglePoolWithAddress(con.Address)
	}
	if con.Pool == nil && con.Kind == KindSentinel {
		con.Pool = pool.NewSentinelPoolWithAddress(con.Address)
	}

	var err error

	var bac backup.Interface
	{
		c := backup.Config{
			Pool: con.Pool,
		}

		bac, err = backup.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var loc locker.Interface
	{
		c := locker.Config{
			Budget: con.Locker.Budget,
			Expiry: con.Locker.Expiry,
			Name:   con.Locker.Name,
			Pool:   con.Pool,
			Prefix: con.Prefix,
		}

		loc, err = locker.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var pub pubsub.Interface
	{
		c := pubsub.Config{
			Pool: con.Pool,

			Prefix: con.Prefix,
		}

		pub, err = pubsub.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var sim simple.Interface
	{
		c := simple.Config{
			Pool: con.Pool,

			Prefix: con.Prefix,
		}

		sim, err = simple.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var sor sorted.Interface
	{
		c := sorted.Config{
			Pool: con.Pool,

			Prefix: con.Prefix,
		}

		sor, err = sorted.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var wal walker.Interface
	{
		c := walker.Config{
			Pool: con.Pool,

			Count:  con.Count,
			Prefix: con.Prefix,
		}

		wal, err = walker.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	c := &Client{
		bac: bac,
		loc: loc,
		poo: con.Pool,
		pub: pub,
		sim: sim,
		sor: sor,
		wal: wal,
	}

	return c, nil
}

func (c *Client) Check() error {
	con := c.poo.Get()
	defer con.Close()

	_, err := con.Do("PING")
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (c *Client) Close() error {
	err := c.poo.Close()
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (c *Client) Empty() (bool, error) {
	con := c.poo.Get()
	defer con.Close()

	res, err := redis.Strings(con.Do("KEYS", "*"))
	if err != nil {
		return false, tracer.Mask(err)
	}

	if len(res) == 0 {
		return true, nil
	}

	return false, nil
}

func (c *Client) Purge() error {
	con := c.poo.Get()
	defer con.Close()

	_, err := con.Do("FLUSHALL")
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (c *Client) Redis(fun func(con redis.Conn) error) error {
	con := c.poo.Get()
	defer con.Close()

	err := fun(con)
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (c *Client) Backup() backup.Interface {
	return c.bac
}

func (c *Client) Locker() locker.Interface {
	return c.loc
}

func (c *Client) PubSub() pubsub.Interface {
	return c.pub
}

func (c *Client) Simple() simple.Interface {
	return c.sim
}

func (c *Client) Sorted() sorted.Interface {
	return c.sor
}

func (c *Client) Walker() walker.Interface {
	return c.wal
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
