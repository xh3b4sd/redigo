package sorted

import (
	"strconv"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"

	"github.com/xh3b4sd/redigo/pkg/pool"
)

func Test_Search_Order_Input_Error(t *testing.T) {
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
			var sea *Search
			{
				con := redigomock.NewConn()
				sea = mustNewSearchWithConn(con)
			}

			_, err := sea.Order("foo", tc.lef, tc.rig)
			if !IsExecutionFailedError(err) {
				t.Fatal("expected", executionFailedError, "got", err)
			}
		})
	}
}

func Test_Search_Order_Input_Valid(t *testing.T) {
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
			var sea *Search
			{
				con := redigomock.NewConn()
				con.Command("ZREVRANGE", "prefix:foo", redigomock.NewAnyInt(), redigomock.NewAnyInt()).Expect([]interface{}{
					[]uint8("one"), []uint8("two"),
				})

				sea = mustNewSearchWithConn(con)
			}

			_, err := sea.Order("foo", tc.lef, tc.rig)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func mustNewSearchWithConn(con redis.Conn) *Search {
	var p *redis.Pool
	{
		p = pool.NewSinglePoolWithConnection(con)
	}

	var s *Search
	{
		s = &Search{
			pool: p,

			prefix: "prefix",
		}
	}

	return s
}
