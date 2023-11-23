//go:build single

package conformance

import (
	"testing"

	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/simple"
)

func Test_Client_Single_Simple_Lifecycle(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := redigo.Config{
			Kind: redigo.KindSingle,
		}

		cli, err = redigo.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		cou, err := cli.Simple().Exists().Multi("foo")
		if err != nil {
			t.Fatal(err)
		}
		if cou != 0 {
			t.Fatalf("element must not exist")
		}
	}

	{
		_, err := cli.Simple().Search().Multi("foo")
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
		cou, err := cli.Simple().Exists().Multi("foo")
		if err != nil {
			t.Fatal(err)
		}
		if cou != 1 {
			t.Fatalf("element must exist")
		}
	}

	{
		res, err := cli.Simple().Search().Multi("foo")
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 1 {
			t.Fatal("expected", 1, "got", len(res))
		}
		if res[0] != "bar" {
			t.Fatal("expected", "bar", "got", res[0])
		}
	}

	{
		res, err := cli.Simple().Search().Multi("foo", "baz")
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 2 {
			t.Fatal("expected", 2, "got", len(res))
		}
		if res[0] != "bar" {
			t.Fatal("expected", "bar", "got", res[0])
		}
		if res[1] != "" {
			t.Fatal("expected", "empty string", "got", res[1])
		}
	}

	{
		err := cli.Simple().Create().Element("baz", "zap")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		res, err := cli.Simple().Search().Multi("foo", "baz")
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 2 {
			t.Fatal("expected", 2, "got", len(res))
		}
		if res[0] != "bar" {
			t.Fatal("expected", "bar", "got", res[0])
		}
		if res[1] != "zap" {
			t.Fatal("expected", "zap", "got", res[1])
		}
	}

	{
		cou, err := cli.Simple().Delete().Multi("foo")
		if err != nil {
			t.Fatal(err)
		}
		if cou != 1 {
			t.Fatalf("element must be deleted")
		}
	}

	{
		cou, err := cli.Simple().Exists().Multi("foo")
		if err != nil {
			t.Fatal(err)
		}
		if cou != 0 {
			t.Fatalf("element must not exist")
		}
	}

	{
		_, err := cli.Simple().Search().Multi("foo")
		if !simple.IsNotFound(err) {
			t.Fatalf("element must not exist")
		}
	}
}
