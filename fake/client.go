package fake

import (
	"github.com/xh3b4sd/redigo"
)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) Ping() error {
	return nil
}

func (c *Client) Scored() redigo.Scored {
	return &Scored{}
}

func (c *Client) Shutdown() {
}

func (c *Client) Simple() redigo.Simple {
	return &Simple{}
}
