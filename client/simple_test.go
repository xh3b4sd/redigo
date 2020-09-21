package client

import (
	"errors"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
)

func Test_Client_Simple_Create_Error(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("SET", "prefix:foo", "bar").ExpectError(executionFailedError)

	c := mustNewClientWithConn(conn)

	err := c.Simple().Create("foo", "bar")
	if !errors.Is(err, executionFailedError) {
		t.Fatal("expected", true, "got", false)
	}
}

func Test_Client_Simple_Create_Success(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("SET", "prefix:foo", "bar").Expect("OK")

	c := mustNewClientWithConn(conn)

	err := c.Simple().Create("foo", "bar")
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
}

func Test_Client_Simple_Delete_Error(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("DEL", "prefix:foo").Expect(int64(0))

	c := mustNewClientWithConn(conn)

	err := c.Simple().Delete("foo")
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
}

func Test_Client_Simple_Delete_Success(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("DEL", "prefix:foo").Expect(int64(1))

	c := mustNewClientWithConn(conn)

	err := c.Simple().Delete("foo")
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
}

func Test_Client_Simple_Exists_False(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("EXISTS", "prefix:foo").Expect(int64(0))

	c := mustNewClientWithConn(conn)

	ok, err := c.Simple().Exists("foo")
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	if ok {
		t.Fatal("expected", false, "got", true)
	}
}

func Test_Client_Simple_Exists_True(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("EXISTS", "prefix:foo").Expect(int64(1))

	c := mustNewClientWithConn(conn)

	ok, err := c.Simple().Exists("foo")
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	if !ok {
		t.Fatal("expected", true, "got", false)
	}
}

func Test_Client_Simple_Search_Error(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("GET", "prefix:foo").ExpectError(executionFailedError)

	c := mustNewClientWithConn(conn)

	_, err := c.Simple().Search("foo")
	if !errors.Is(err, executionFailedError) {
		t.Fatal("expected", true, "got", false)
	}
}

func Test_Client_Simple_Search_Error_NotFound(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("GET", "prefix:foo").ExpectError(redis.ErrNil)

	c := mustNewClientWithConn(conn)

	_, err := c.Simple().Search("foo")
	if !IsNotFound(err) {
		t.Fatal("expected", true, "got", false)
	}
}

func Test_Client_Simple_Search_Success(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("GET", "prefix:foo").Expect("bar")

	c := mustNewClientWithConn(conn)

	value, err := c.Simple().Search("foo")
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	if value != "bar" {
		t.Fatal("expected", "bar", "got", value)
	}
}
