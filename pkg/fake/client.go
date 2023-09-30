package fake

import (
	"github.com/gomodule/redigo/redis"

	"github.com/xh3b4sd/redigo/pkg/backup"
	"github.com/xh3b4sd/redigo/pkg/locker"
	"github.com/xh3b4sd/redigo/pkg/pubsub"
	"github.com/xh3b4sd/redigo/pkg/simple"
	"github.com/xh3b4sd/redigo/pkg/sorted"
	"github.com/xh3b4sd/redigo/pkg/walker"
)

type Client struct {
	FakeBackup func() backup.Interface
	FakeLocker func() locker.Interface
	FakePubSub func() pubsub.Interface
	FakeSimple func() simple.Interface
	FakeSorted func() sorted.Interface
	FakeWalker func() walker.Interface
}

func New() *Client {
	return &Client{}
}

func (c *Client) Check() error {
	return nil
}

func (c *Client) Close() error {
	return nil
}

func (c *Client) Empty() (bool, error) {
	return false, nil
}

func (c *Client) Purge() error {
	return nil
}

func (c *Client) Redis(fun func(con redis.Conn) error) error {
	return nil
}

func (c *Client) Listen() string {
	return ""
}

func (c *Client) Backup() backup.Interface {
	if c.FakeLocker != nil {
		return c.FakeBackup()
	}

	return &Backup{}
}

func (c *Client) Locker() locker.Interface {
	if c.FakeLocker != nil {
		return c.FakeLocker()
	}

	return &Locker{}
}

func (c *Client) PubSub() pubsub.Interface {
	if c.FakePubSub != nil {
		return c.FakePubSub()
	}

	return &PubSub{}
}

func (c *Client) Simple() simple.Interface {
	if c.FakeSimple != nil {
		return c.FakeSimple()
	}

	return &Simple{}
}

func (c *Client) Sorted() sorted.Interface {
	if c.FakeSorted != nil {
		return c.FakeSorted()
	}

	return &Sorted{}
}

func (c *Client) Walker() walker.Interface {
	if c.FakeWalker != nil {
		return c.FakeWalker()
	}

	return &Walker{}
}
