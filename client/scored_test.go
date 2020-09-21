package client

import (
	"errors"
	"testing"

	"github.com/rafaeljusto/redigomock"
)

func Test_Client_Scored_Create_Success(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("ZADD", "prefix:key", 0.8, "element").Expect(int64(1))

	c := mustNewClientWithConn(conn)

	err := c.Scored().Create("key", "element", 0.8)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
}

func Test_Client_Scored_Create_Error(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("ZADD", "prefix:key", 0.8, "element").ExpectError(executionFailedError)

	c := mustNewClientWithConn(conn)

	err := c.Scored().Create("key", "element", 0.8)
	if !errors.Is(err, executionFailedError) {
		t.Fatal("expected", true, "got", false)
	}
}

func Test_Client_Scored_Delete_Success(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("ZREM", "prefix:test-key", "test-element").Expect(int64(1))

	c := mustNewClientWithConn(conn)

	err := c.Scored().Delete("test-key", "test-element")
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
}

func Test_Client_Scored_Delete_Error(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("ZREM", "prefix:test-key", "test-element").ExpectError(executionFailedError)

	c := mustNewClientWithConn(conn)

	err := c.Scored().Delete("test-key", "test-element")
	if !errors.Is(err, executionFailedError) {
		t.Fatal("expected", true, "got", false)
	}
}

// Test_Client_Scored_CutOff_Success_CutOff tests a case where the scored set
// tried to be truncated is big enough to be shortened. This is because the
// maximum length given to CutOff is 10 while the ZCARD command here returns 12,
// meaning 2 elements have to be removed from the sorted set.
func Test_Client_Scored_CutOff_Success_CutOff(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("ZCARD", "prefix:foo").Expect(int64(12))
	conn.Command("ZPOPMIN", "prefix:foo", 2).Expect([]interface{}{
		[]uint8("25"), []uint8("one"), []uint8("35"), []uint8("two"),
	})

	c := mustNewClientWithConn(conn)

	err := c.Scored().CutOff("foo", 10)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
}

// Test_Client_Scored_CutOff_Success_NoCutOff tests a case where the scored set
// tried to be truncated is not big enough to be shortened. This is because the
// maximum length given to CutOff is 10 while the ZCARD command here only
// returns 3, meaning nothing has to be done.
func Test_Client_Scored_CutOff_Success_NoCutOff(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("ZCARD", "prefix:foo").Expect(int64(3))

	c := mustNewClientWithConn(conn)

	err := c.Scored().CutOff("foo", 10)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
}

func Test_Client_Scored_Search_Success(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("ZREVRANGE", "prefix:foo", 0, 1, "WITHSCORES").Expect([]interface{}{
		[]uint8("one"), []uint8("0.8"), []uint8("two"), []uint8("0.5"),
	})

	c := mustNewClientWithConn(conn)

	values, err := c.Scored().Search("foo", 2)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	if len(values) != 4 {
		t.Fatal("expected", 1, "got", len(values))
	}
	if values[0] != "one" {
		t.Fatal("expected", "one", "got", values[0])
	}
	if values[1] != "0.8" {
		t.Fatal("expected", "0.8", "got", values[1])
	}
	if values[2] != "two" {
		t.Fatal("expected", "two", "got", values[2])
	}
	if values[3] != "0.5" {
		t.Fatal("expected", "0.5", "got", values[3])
	}
}

func Test_Client_Scored_Search_Error(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("ZREVRANGE", "prefix:foo", 0, 1, "WITHSCORES").ExpectError(executionFailedError)

	c := mustNewClientWithConn(conn)

	_, err := c.Scored().Search("foo", 2)
	if !errors.Is(err, executionFailedError) {
		t.Fatal("expected", true, "got", false)
	}
}
