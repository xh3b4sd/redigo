package client

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/rafaeljusto/redigomock"
	"github.com/xh3b4sd/tracer"
)

func Test_Client_Scored_Create_Success(t *testing.T) {
	con := redigomock.NewConn()
	con.Command("ZADD", "prefix:key", 0.8, "element").Expect(int64(1))

	cli := mustNewClientWithConn(con)

	err := cli.Scored().Create("key", "element", 0.8)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Client_Scored_Create_Error(t *testing.T) {
	con := redigomock.NewConn()
	con.Command("ZADD", "prefix:key", 0.8, "element").ExpectError(executionFailedError)

	cli := mustNewClientWithConn(con)

	err := cli.Scored().Create("key", "element", 0.8)
	if !errors.Is(err, executionFailedError) {
		t.Fatal("expected", true, "got", false)
	}
}

func Test_Client_Scored_Delete_Success(t *testing.T) {
	con := redigomock.NewConn()
	con.Command("ZREM", "prefix:test-key", "test-element").Expect(int64(1))

	cli := mustNewClientWithConn(con)

	err := cli.Scored().Delete("test-key", "test-element")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Client_Scored_Delete_Error(t *testing.T) {
	con := redigomock.NewConn()
	con.Command("ZREM", "prefix:test-key", "test-element").ExpectError(executionFailedError)

	cli := mustNewClientWithConn(con)

	err := cli.Scored().Delete("test-key", "test-element")
	if !errors.Is(err, executionFailedError) {
		t.Fatal("expected", true, "got", false)
	}
}

func Test_Client_Scored_Exists_Error(t *testing.T) {
	testError := &tracer.Error{Kind: "testError"}

	testCases := []struct {
		err error
	}{
		// Case 0 ensures that redis errors are populated properly.
		{
			err: testError,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var cli *Client
			{
				con := redigomock.NewConn()
				con.Command("EXISTS", "prefix:foo").ExpectError(tc.err)

				cli = mustNewClientWithConn(con)
			}

			_, err := cli.Scored().Exists("foo")
			if !errors.Is(err, tc.err) {
				t.Fatal("expected error to match")
			}
		})
	}
}

func Test_Client_Scored_Exists_Input(t *testing.T) {
	testCases := []struct {
		k string
		i int
		b bool
	}{
		// Case 0 ensures that redis 0 means exists false.
		{
			k: "foo",
			i: 0,
			b: false,
		},
		// Case 1 ensures that redis 0 means exists false.
		{
			k: "bar",
			i: 0,
			b: false,
		},
		// Case 2 ensures that redis 1 means exists true.
		{
			k: "baz",
			i: 1,
			b: true,
		},
		// Case 3 ensures that redis 1 means exists true.
		{
			k: "zap",
			i: 1,
			b: true,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var cli *Client
			{
				con := redigomock.NewConn()
				con.Command("EXISTS", fmt.Sprintf("prefix:%s", tc.k)).Expect(int64(tc.i))

				cli = mustNewClientWithConn(con)
			}

			ok, err := cli.Scored().Exists(tc.k)
			if err != nil {
				t.Fatal(err)
			}
			if ok != tc.b {
				t.Fatal("expected", nil, "got", err)
			}
		})
	}
}

func Test_Client_Scored_Search_Input_Error(t *testing.T) {
	testCases := []struct {
		key string
		lef int
		rig int
	}{
		// Case 0 ensures that lef cannot be negative.
		{
			key: "foo",
			lef: -1,
			rig: 1,
		},
		// Case 1 ensures that lef cannot be negative.
		{
			key: "bar",
			lef: -6,
			rig: 1,
		},
		// Case 2 ensures that rig cannot be smaller than -1.
		{
			key: "baz",
			lef: 0,
			rig: -2,
		},
		// Case 3 ensures that rig cannot be smaller than -1.
		{
			key: "zap",
			lef: 0,
			rig: -4,
		},
		// Case 4 ensures that lef cannot be greater than rig.
		{
			key: "ler",
			lef: 3,
			rig: 2,
		},
		// Case 5 ensures that lef cannot be greater than rig.
		{
			key: "tml",
			lef: 10,
			rig: 5,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var cli *Client
			{
				con := redigomock.NewConn()
				con.Command("ZREVRANGE", fmt.Sprintf("prefix:%s", tc.key), redigomock.NewAnyInt(), redigomock.NewAnyInt()).Expect([]interface{}{
					[]uint8("one"), []uint8("two"),
				})

				cli = mustNewClientWithConn(con)
			}

			_, err := cli.Scored().Search(tc.key, tc.lef, tc.rig)
			if !errors.Is(err, executionFailedError) {
				t.Fatal("expected", nil, "got", err)
			}
		})
	}
}

func Test_Client_Scored_Search_Input_Valid(t *testing.T) {
	testCases := []struct {
		key string
		lef int
		rig int
	}{
		// Case 0 ensures that a single element can be searched.
		{
			key: "foo",
			lef: 0,
			rig: 1,
		},
		// Case 1 ensures that all elements can be searched.
		{
			key: "bar",
			lef: 0,
			rig: -1,
		},
		// Case 2 ensures that multiple elements can be searched.
		{
			key: "baz",
			lef: 0,
			rig: 3,
		},
		// Case 3 ensures that a single element can be searched within the
		// dataset.
		{
			key: "zap",
			lef: 4,
			rig: 5,
		},
		// Case 4 ensures that multiple elements can be searched within the
		// dataset.
		{
			key: "ler",
			lef: 10,
			rig: 20,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var cli *Client
			{
				con := redigomock.NewConn()
				con.Command("ZREVRANGE", fmt.Sprintf("prefix:%s", tc.key), redigomock.NewAnyInt(), redigomock.NewAnyInt()).Expect([]interface{}{
					[]uint8("one"), []uint8("two"),
				})

				cli = mustNewClientWithConn(con)
			}

			_, err := cli.Scored().Search(tc.key, tc.lef, tc.rig)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func Test_Client_Scored_Search_Redis_Error(t *testing.T) {
	testError := &tracer.Error{Kind: "testError"}

	testCases := []struct {
		err error
	}{
		// Case 0 ensures that redis errors are populated properly.
		{
			err: testError,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var cli *Client
			{
				con := redigomock.NewConn()
				con.Command("ZREVRANGE", "prefix:foo", 0, 1).ExpectError(tc.err)

				cli = mustNewClientWithConn(con)
			}

			_, err := cli.Scored().Search("foo", 0, 2)
			if !errors.Is(err, tc.err) {
				t.Fatal("expected", true, "got", false)
			}
		})
	}
}

func Test_Client_Scored_Search_Redis_Valid(t *testing.T) {
	con := redigomock.NewConn()
	con.Command("ZREVRANGE", "prefix:foo", 0, 1).Expect([]interface{}{
		[]uint8("one"), []uint8("two"),
	})

	cli := mustNewClientWithConn(con)

	values, err := cli.Scored().Search("foo", 0, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(values) != 2 {
		t.Fatal("expected", 2, "got", len(values))
	}
	if values[0] != "one" {
		t.Fatal("expected", "one", "got", values[0])
	}
	if values[1] != "two" {
		t.Fatal("expected", "two", "got", values[1])
	}
}
