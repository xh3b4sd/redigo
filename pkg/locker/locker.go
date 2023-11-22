package locker

import (
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/redigo"
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/breakr"
	"github.com/xh3b4sd/redigo/pkg/prefix"
)

type Config struct {
	Breakr breakr.Interface
	Expiry time.Duration
	Name   string
	Pool   *redis.Pool
	Prefix string
}

type Locker struct {
	brk breakr.Interface
	mut *redsync.Mutex
}

func New(c Config) *Locker {
	if c.Breakr == nil {
		c.Breakr = breakr.Default()
	}
	if c.Expiry == 0 {
		c.Expiry = 30 * time.Second
	}
	if c.Name == "" {
		c.Name = "def"
	}

	var r *redsync.Redsync
	{
		p := redigo.NewPool(
			&redis.Pool{
				MaxIdle:      c.Pool.MaxIdle,
				IdleTimeout:  c.Pool.IdleTimeout,
				Dial:         c.Pool.Dial,
				TestOnBorrow: c.Pool.TestOnBorrow,
			},
		)

		r = redsync.New(p)
	}

	var l *Locker
	{
		l = &Locker{
			brk: c.Breakr,
			mut: r.NewMutex(
				prefix.WithKeys(c.Prefix, fmt.Sprintf("red:%s", c.Name)),
				redsync.WithExpiry(c.Expiry),
				redsync.WithRetryDelayFunc(func(tries int) time.Duration { return 100 * time.Millisecond }),
				redsync.WithTries(1),
			),
		}
	}

	return l
}
