// +build redis

package client

import (
	"testing"

	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/sorted"
)

func Test_Client_Sorted_Redis_Exists(t *testing.T) {
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
		exi, err := cli.Sorted().Exists().Score("ssk", 0.8)
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatalf("element must not exist")
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

	{
		err := cli.Sorted().Create().Element("ssk", "foo", 0.8)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Score("ssk", 0.8)
		if err != nil {
			t.Fatal(err)
		}
		if !exi {
			t.Fatalf("element must exist")
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
			t.Fatal(err)
		}
	}

	{
		exi, err := cli.Sorted().Exists().Score("ssk", 0.8)
		if err != nil {
			t.Fatal(err)
		}
		if exi {
			t.Fatalf("element must not exist")
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

func Test_Client_Sorted_Redis_Index(t *testing.T) {
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
		err := cli.Sorted().Create().Element("ssk", "foo", 0.8, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "foo", 0.7, "c", "d")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "foo", 0.6, "c", "d")
		if !sorted.IsAlreadyExistsError(err) {
			t.Fatal("expected", alreadyExistsError, "got", err)
		}
	}

	{
		err := cli.Sorted().Create().Element("ssk", "bar", 0.8, "e", "f")
		if !sorted.IsAlreadyExistsError(err) {
			t.Fatal("expected", alreadyExistsError, "got", err)
		}
	}
}
