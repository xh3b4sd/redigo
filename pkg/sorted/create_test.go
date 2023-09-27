package sorted

import (
	"strconv"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"

	"github.com/xh3b4sd/redigo/pkg/pool"
)

func Test_Create_Score_Input_Error(t *testing.T) {
	testCases := []struct {
		ind []string
	}{
		// Case 0 ensures that indices must not be empty.
		{
			ind: []string{
				"",
			},
		},
		// Case 1 ensures that indices must not be empty.
		{
			ind: []string{
				"",
				"",
			},
		},
		// Case 2 ensures that indices must not contain whitespace.
		{
			ind: []string{
				" ",
			},
		},
		// Case 3 ensures that indices must not contain whitespace.
		{
			ind: []string{
				"  ",
				" foo ",
			},
		},
		// Case 4 ensures that indices must not be duplicated.
		{
			ind: []string{
				"a",
				"a",
			},
		},
		// Case 5 ensures that indices must not be duplicated.
		{
			ind: []string{
				"b:2",
				"b:2",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var cre *create
			{
				con := redigomock.NewConn()
				cre = mustNewCreateWithConn(con)
			}

			err := cre.Score("ssk", "foo", 0.8, tc.ind...)
			if !IsExecutionFailedError(err) {
				t.Fatal("expected", executionFailedError, "got", err)
			}
		})
	}
}

func Test_Create_Score_Input_Valid(t *testing.T) {
	testCases := []struct {
		ind []string
	}{
		// Case 0 ensures that no indices are valid.
		{
			ind: []string{},
		},
		// Case 1 ensures that indices are valid.
		{
			ind: []string{
				"a:1",
			},
		},
		// Case 2 ensures that indices are valid.
		{
			ind: []string{
				"a:1",
				"b:2",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var cre *create
			{
				con := redigomock.NewConn()
				con.GenericCommand("EVALSHA").Expect(int64(2))

				cre = mustNewCreateWithConn(con)
			}

			err := cre.Score("ssk", "foo", 0.8, tc.ind...)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func mustNewCreateWithConn(con redis.Conn) *create {
	var p *redis.Pool
	{
		p = pool.NewSinglePoolWithConnection(con)
	}

	var c *create
	{
		c = &create{
			pool: p,

			createScoreScript: redis.NewScript(2, createScoreScript),

			prefix: "prefix",
		}
	}

	return c
}
