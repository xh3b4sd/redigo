package sorted

import (
	"errors"
	"strconv"
	"testing"

	"github.com/rafaeljusto/redigomock"
	"github.com/xh3b4sd/tracer"
)

func Test_Client_Scored_Create_Error(t *testing.T) {
	con := redigomock.NewConn()
	con.Command("EVALSHA", redigomock.NewAnyData(), 1, "prefix:key", "ele", 0.8).Expect(int64(0))

	cli := mustNewClientWithConn(con)

	err := cli.Scored().Create("key", "ele", 0.8)
	if !errors.Is(err, alreadyExistsError) {
		t.Fatal("expected", true, "got", false)
	}
}

func Test_Client_Scored_Create_Success(t *testing.T) {
	con := redigomock.NewConn()
	con.Command("EVALSHA", redigomock.NewAnyData(), 1, "prefix:key", "ele", 0.8).Expect(int64(1))

	cli := mustNewClientWithConn(con)

	err := cli.Scored().Create("key", "ele", 0.8)
	if err != nil {
		t.Fatal(err)
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
				con.Command("ZRANGEBYSCORE", "prefix:foo", 0.8, 0.8).ExpectError(tc.err)

				cli = mustNewClientWithConn(con)
			}

			_, err := cli.Scored().Exists("foo", 0.8)
			if !errors.Is(err, tc.err) {
				t.Fatal("expected error to match")
			}
		})
	}
}

func Test_Client_Scored_Exists_Input(t *testing.T) {
	testCases := []struct {
		res []interface{}
		exi bool
	}{
		// Case 0 ensures that elements do not exist.
		{
			res: []interface{}{},
			exi: false,
		},
		// Case 1 ensures that elements exist.
		{
			res: []interface{}{
				[]uint8("one"),
			},
			exi: true,
		},
		// Case 2 ensures that elements exist.
		{
			res: []interface{}{
				[]uint8("one"),
				[]uint8("two"),
			},
			exi: true,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var cli *Client
			{
				con := redigomock.NewConn()
				con.Command("ZRANGEBYSCORE", "prefix:foo", 0.8, 0.8).Expect(tc.res)

				cli = mustNewClientWithConn(con)
			}

			exi, err := cli.Scored().Exists("foo", 0.8)
			if err != nil {
				t.Fatal(err)
			}
			if exi != tc.exi {
				t.Fatal("expected", tc.exi, "got", exi)
			}
		})
	}
}

func Test_Client_Scored_Search_Input_Error(t *testing.T) {
	testCases := []struct {
		lef int
		rig int
	}{
		// Case 0 ensures that lef cannot be negative.
		{
			lef: -1,
			rig: 1,
		},
		// Case 1 ensures that lef cannot be negative.
		{
			lef: -6,
			rig: 1,
		},
		// Case 2 ensures that rig cannot be smaller than -1.
		{
			lef: 0,
			rig: -2,
		},
		// Case 3 ensures that rig cannot be smaller than -1.
		{
			lef: 0,
			rig: -4,
		},
		// Case 4 ensures that lef cannot be greater than rig.
		{
			lef: 3,
			rig: 2,
		},
		// Case 5 ensures that lef cannot be greater than rig.
		{
			lef: 10,
			rig: 5,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var cli *Client
			{
				con := redigomock.NewConn()
				con.Command("ZREVRANGE", "prefix:foo", redigomock.NewAnyInt(), redigomock.NewAnyInt()).Expect([]interface{}{
					[]uint8("one"), []uint8("two"),
				})

				cli = mustNewClientWithConn(con)
			}

			_, err := cli.Scored().Search("foo", tc.lef, tc.rig)
			if !errors.Is(err, executionFailedError) {
				t.Fatal("expected", nil, "got", err)
			}
		})
	}
}

func Test_Client_Scored_Search_Input_Valid(t *testing.T) {
	testCases := []struct {
		lef int
		rig int
	}{
		// Case 0 ensures that a single element can be searched.
		{
			lef: 0,
			rig: 1,
		},
		// Case 1 ensures that all elements can be searched.
		{
			lef: 0,
			rig: -1,
		},
		// Case 2 ensures that multiple elements can be searched.
		{
			lef: 0,
			rig: 3,
		},
		// Case 3 ensures that a single element can be searched within the
		// dataset.
		{
			lef: 4,
			rig: 5,
		},
		// Case 4 ensures that multiple elements can be searched within the
		// dataset.
		{
			lef: 10,
			rig: 20,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var cli *Client
			{
				con := redigomock.NewConn()
				con.Command("ZREVRANGE", "prefix:foo", redigomock.NewAnyInt(), redigomock.NewAnyInt()).Expect([]interface{}{
					[]uint8("one"), []uint8("two"),
				})

				cli = mustNewClientWithConn(con)
			}

			_, err := cli.Scored().Search("foo", tc.lef, tc.rig)
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

func Test_Client_Scored_Update_Call(t *testing.T) {
	con := redigomock.NewConn()
	con.Command("EVALSHA", redigomock.NewAnyData(), 1, "prefix:key", "ele", 0.8).Expect(int64(3))

	cli := mustNewClientWithConn(con)

	upd, err := cli.Scored().Update("key", "ele", 0.8)
	if err != nil {
		t.Fatal(err)
	}
	if !upd {
		t.Fatal(err)
	}
}
