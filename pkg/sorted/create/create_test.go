package create

import (
	"errors"
	"strconv"
	"testing"

	"github.com/rafaeljusto/redigomock"
	"github.com/xh3b4sd/redigo/pkg/pool"
)

func Test_Sorted_Create_Index_Input_Error(t *testing.T) {
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
			var r *Redis
			{
				r = New(Config{
					Poo: pool.NewSinglePoolWithConnection(redigomock.NewConn()),
				})
			}

			err := r.Index("ssk", "foo", 0.8, tc.ind...)
			if !errors.Is(err, executionFailedError) {
				t.Fatal("expected", executionFailedError, "got", err)
			}
		})
	}
}
