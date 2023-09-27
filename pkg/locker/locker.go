package locker

import (
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/redigo"
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/budget/v3"
	"github.com/xh3b4sd/budget/v3/pkg/breaker"

	"github.com/xh3b4sd/redigo/pkg/prefix"
)

type Config struct {
	Budget budget.Interface
	Expiry time.Duration
	Name   string
	Pool   *redis.Pool
	Prefix string
}

type Locker struct {
	bud   budget.Interface
	mutex *redsync.Mutex
}

func New(config Config) (*Locker, error) {
	if config.Budget == nil {
		config.Budget = breaker.Default()
	}
	if config.Expiry == 0 {
		config.Expiry = 30 * time.Second
	}
	if config.Name == "" {
		config.Name = "def"
	}

	var r *redsync.Redsync
	{
		p := redigo.NewPool(
			&redis.Pool{
				MaxIdle:      config.Pool.MaxIdle,
				IdleTimeout:  config.Pool.IdleTimeout,
				Dial:         config.Pool.Dial,
				TestOnBorrow: config.Pool.TestOnBorrow,
			},
		)

		r = redsync.New(p)
	}

	l := &Locker{
		bud: config.Budget,
		mutex: r.NewMutex(
			prefix.WithKeys(config.Prefix, fmt.Sprintf("red:%s", config.Name)),
			redsync.WithExpiry(config.Expiry),
			redsync.WithRetryDelayFunc(func(tries int) time.Duration { return 100 * time.Millisecond }),
			redsync.WithTries(1),
		),
	}

	return l, nil
}
