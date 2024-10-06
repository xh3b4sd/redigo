package redigo

import (
	"net"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/redigo/pkg/backup"
	"github.com/xh3b4sd/redigo/pkg/pubsub"
	"github.com/xh3b4sd/redigo/pkg/sorted"
	"github.com/xh3b4sd/redigo/pool"
	"github.com/xh3b4sd/redigo/simple"
	"github.com/xh3b4sd/redigo/walker"
	"github.com/xh3b4sd/tracer"
)

const (
	KindSingle   = "single"
	KindSentinel = "sentinel"
)

type Config struct {
	Address string
	Count   int
	Kind    string
	Pass    string
	Pool    *redis.Pool
	Prefix  string
	User    string
}

type Redigo struct {
	add string
	bac backup.Interface
	poo *redis.Pool
	pub pubsub.Interface
	sim simple.Interface
	sor sorted.Interface
	wal walker.Interface
}

func New(con Config) (*Redigo, error) {
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
		con.Pool = pool.NewSinglePoolWithAddress(con.Address, con.User, con.Pass)
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
		sim = simple.New(simple.Config{
			Pool: con.Pool,

			Prefix: con.Prefix,
		})
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
		wal = walker.New(walker.Config{
			Pool: con.Pool,

			Count:  con.Count,
			Prefix: con.Prefix,
		})
	}

	r := &Redigo{
		add: con.Address,
		bac: bac,
		poo: con.Pool,
		pub: pub,
		sim: sim,
		sor: sor,
		wal: wal,
	}

	return r, nil
}

func (r *Redigo) Check() error {
	con := r.poo.Get()
	defer con.Close()

	_, err := con.Do("PING")
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (r *Redigo) Close() error {
	err := r.poo.Close()
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (r *Redigo) Empty() (bool, error) {
	con := r.poo.Get()
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

func (r *Redigo) Purge() error {
	con := r.poo.Get()
	defer con.Close()

	_, err := con.Do("FLUSHALL")
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (r *Redigo) Redis(fun func(con redis.Conn) error) error {
	con := r.poo.Get()
	defer con.Close()

	err := fun(con)
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (r *Redigo) Listen() string {
	return r.add
}

func (r *Redigo) Backup() backup.Interface {
	return r.bac
}

func (r *Redigo) PubSub() pubsub.Interface {
	return r.pub
}

func (r *Redigo) Simple() simple.Interface {
	return r.sim
}

func (r *Redigo) Sorted() sorted.Interface {
	return r.sor
}

func (r *Redigo) Walker() walker.Interface {
	return r.wal
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
