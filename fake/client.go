package fake

import (
	"github.com/xh3b4sd/redigo"
)

type Client struct {
	PubSubFake func() redigo.PubSub
	SortedFake func() redigo.Sorted
	SimpleFake func() redigo.Simple
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

func (c *Client) PubSub() redigo.PubSub {
	if c.PubSubFake != nil {
		return c.PubSubFake()
	}

	return &PubSub{}
}

func (c *Client) Sorted() redigo.Sorted {
	if c.SortedFake != nil {
		return c.SortedFake()
	}

	return &Sorted{}
}

func (c *Client) Simple() redigo.Simple {
	if c.SimpleFake != nil {
		return c.SimpleFake()
	}

	return &Simple{}
}
