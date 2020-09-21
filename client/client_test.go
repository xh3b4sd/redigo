package client

import (
	"testing"

	"github.com/gomodule/redigo/redis"
)

func Test_Client_Shutdown(t *testing.T) {
	c := mustNewClientWithConn(nil)

	c.Shutdown()
	c.Shutdown()
	c.Shutdown()
}

func Test_Client_withPrefix(t *testing.T) {
	expected := "test-prefix:my:test:key"
	newKey := withPrefix("test-prefix", "my", "test", "key")
	if newKey != expected {
		t.Fatal("expected", expected, "got", newKey)
	}
}

func Test_Client_withPrefix_Empty(t *testing.T) {
	newKey := withPrefix("test-prefix")
	if newKey != "test-prefix" {
		t.Fatal("expected", "test-prefix", "got", newKey)
	}
}

func mustNewClientWithConn(conn redis.Conn) *Client {
	var err error

	var pool *redis.Pool
	{
		c := DefaultPoolConfig()
		c.Dial = func() (redis.Conn, error) { return conn, nil }

		pool = NewPool(c)
	}

	var client *Client
	{
		c := Config{
			Pool:   pool,
			Prefix: "prefix",
		}

		client, err = New(c)
		if err != nil {
			panic(err)
		}
	}

	return client
}
