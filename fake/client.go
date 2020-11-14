package fake

import (
	"github.com/xh3b4sd/redigo"
)

type Client struct {
	ScoredFake func() redigo.Scored
	SimpleFake func() redigo.Simple
}

func New() *Client {
	return &Client{}
}

func (c *Client) Ping() error {
	return nil
}

func (c *Client) Scored() redigo.Scored {
	if c.ScoredFake != nil {
		return c.ScoredFake()
	}

	return &Scored{}
}

func (c *Client) Shutdown() {
}

func (c *Client) Simple() redigo.Simple {
	if c.SimpleFake != nil {
		return c.SimpleFake()
	}

	return &Simple{}
}
