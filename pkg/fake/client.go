package fake

import (
	"github.com/xh3b4sd/redigo"
)

type Client struct {
	FakeLocker func() redigo.Locker
	FakePubSub func() redigo.PubSub
	FakeSorted func() redigo.Sorted
	FakeSimple func() redigo.Simple
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

func (c *Client) Purge() error {
	return nil
}

func (c *Client) Locker() redigo.Locker {
	if c.FakeLocker != nil {
		return c.FakeLocker()
	}

	return &Locker{}
}

func (c *Client) PubSub() redigo.PubSub {
	if c.FakePubSub != nil {
		return c.FakePubSub()
	}

	return &PubSub{}
}

func (c *Client) Sorted() redigo.Sorted {
	if c.FakeSorted != nil {
		return c.FakeSorted()
	}

	return &Sorted{}
}

func (c *Client) Simple() redigo.Simple {
	if c.FakeSimple != nil {
		return c.FakeSimple()
	}

	return &Simple{}
}
