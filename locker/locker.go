package locker

import (
	redsync "github.com/go-redsync/redsync/v4"
	redsyncredigo "github.com/go-redsync/redsync/v4/redis/redigo"
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/redigo/prefix"
)

type Config struct {
	Pool *redis.Pool

	Prefix string
}

type Locker struct {
	mutex *redsync.Mutex
}

func New(config Config) (*Locker, error) {
	var r *redsync.Redsync
	{
		p := redsyncredigo.NewPool(
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
		mutex: r.NewMutex(prefix.WithKeys(config.Prefix, "redigo:locker")),
	}

	return l, nil
}
