package redigo

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/redigo/pkg/backup"
	"github.com/xh3b4sd/redigo/pkg/pubsub"
	"github.com/xh3b4sd/redigo/pkg/sorted"
	"github.com/xh3b4sd/redigo/pkg/walker"
	"github.com/xh3b4sd/redigo/simple"
)

type Interface interface {
	Check() error
	Close() error
	Empty() (bool, error)
	Purge() error
	Redis(fun func(con redis.Conn) error) error

	// Listen returns the host:port configuration for this redigo instance.
	Listen() string

	Backup() backup.Interface
	PubSub() pubsub.Interface
	Sorted() sorted.Interface
	Simple() simple.Interface
	Walker() walker.Interface
}
