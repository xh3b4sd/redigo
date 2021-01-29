// +build sentinel

package client

import (
	"testing"

	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/pkg/simple"
)

func Test_Client_Sentinel_Connection(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := Config{
			Kind: KindSentinel,
		}

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
		exi, err := cli.Simple().Exists().Element("foo")
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatalf("element must not exist")
		}
	}

	{
		_, err := cli.Simple().Search().Value("foo")
		if !simple.IsNotFound(err) {
			t.Fatalf("element must not exist")
		}
	}

	{
		err := cli.Simple().Create().Element("foo", "bar")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Simple().Exists().Element("foo")
		if err != nil {
			t.Fatal(err)
		}
		if !exi {
			t.Fatalf("element must exist")
		}
	}

	{
		val, err := cli.Simple().Search().Value("foo")
		if err != nil {
			t.Fatal(err)
		}
		if val != "bar" {
			t.Fatalf("val must be bar")
		}
	}

	{
		err := cli.Simple().Delete().Element("foo")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Simple().Exists().Element("foo")
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatalf("element must not exist")
		}
	}

	{
		_, err := cli.Simple().Search().Value("foo")
		if !simple.IsNotFound(err) {
			t.Fatalf("element must not exist")
		}
	}
}
