package backup

import (
	"github.com/gomodule/redigo/redis"
)

type Config struct {
	Pool *redis.Pool
}

type Backup struct {
	poo *redis.Pool
}

func New(con Config) (*Backup, error) {
	b := &Backup{
		poo: con.Pool,
	}

	return b, nil
}
