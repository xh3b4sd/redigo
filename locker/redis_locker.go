package locker

import (
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/redigo"
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/breakr"
	"github.com/xh3b4sd/redigo/prefix"
)

type Config struct {
	// Brk is the budget implementation used to retry redis connections on
	// failure.
	Brk breakr.Interface
	// Exp is the lock's expiry, so that locks can expire after a certain amount
	// of time of inactivity. Defaults to 30 seconds. Disabled with -1.
	Exp time.Duration
	// Nam is some uniquely scoped identifier for a lock purpose.
	Nam string
	// Poo is the redis connection pool to select client connections from.
	Poo *redis.Pool
	// Pre is the prefix of the underlying redis key used to coordinate the
	// distributed lock.
	Pre string
}

type Redis struct {
	brk breakr.Interface
	mut *redsync.Mutex
}

func New(c Config) *Redis {
	if c.Brk == nil {
		c.Brk = breakr.Default()
	}
	if c.Exp == 0 {
		c.Exp = 30 * time.Second
	}
	if c.Nam == "" {
		c.Nam = "def"
	}

	var r *redsync.Redsync
	{
		p := redigo.NewPool(
			&redis.Pool{
				MaxIdle:      c.Poo.MaxIdle,
				IdleTimeout:  c.Poo.IdleTimeout,
				Dial:         c.Poo.Dial,
				TestOnBorrow: c.Poo.TestOnBorrow,
			},
		)

		r = redsync.New(p)
	}

	var o []redsync.Option
	{
		o = []redsync.Option{
			redsync.WithRetryDelayFunc(func(tries int) time.Duration { return 100 * time.Millisecond }),
			redsync.WithTries(1),
		}
	}

	if c.Exp > 0 {
		o = append(o, redsync.WithExpiry(c.Exp))
	}

	var l *Redis
	{
		l = &Redis{
			brk: c.Brk,
			mut: r.NewMutex(
				prefix.WithKeys(c.Pre, fmt.Sprintf("red:%s", c.Nam)),
				o...,
			),
		}
	}

	return l
}
