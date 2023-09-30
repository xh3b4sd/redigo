package redigo

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/redigo/pkg/backup"
	"github.com/xh3b4sd/redigo/pkg/locker"
	"github.com/xh3b4sd/redigo/pkg/pubsub"
	"github.com/xh3b4sd/redigo/pkg/simple"
	"github.com/xh3b4sd/redigo/pkg/sorted"
	"github.com/xh3b4sd/redigo/pkg/walker"
)

type Interface interface {
	Check() error
	Close() error
	Empty() (bool, error)
	Purge() error
	Redis(fun func(con redis.Conn) error) error

	Listen() string

	Backup() backup.Interface
	Locker() locker.Interface
	PubSub() pubsub.Interface
	Sorted() sorted.Interface
	Simple() simple.Interface
	Walker() walker.Interface
}
