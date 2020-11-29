// +build redis

package client

import (
	"fmt"
	"testing"

	"github.com/xh3b4sd/redigo"
)

func Test_Client_Sorted_Redis(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := Config{}

		cli, err = New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Value("ssk", "foo")
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatalf("value must not exist")
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "foo", 0.8)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Value("ssk", "foo")
		if err != nil {
			t.Fatal(err)
		}
		if !exi {
			t.Fatalf("element must exist")
		}
	}

	{
		err := cli.Sorted().Delete().Element("ssk", "foo")
		if err != nil {
			fmt.Printf("%#v\n", err)
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Value("ssk", "foo")
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatalf("element must not exist")
		}
	}
}
