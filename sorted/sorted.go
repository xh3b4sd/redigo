package sorted

import (
	"github.com/gomodule/redigo/redis"
)

type Sorted struct {
	pool *redis.Pool

	prefix string
}
