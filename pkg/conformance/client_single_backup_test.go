//go:build single

package conformance

import (
	"testing"
	"time"

	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/pkg/client"
)

func Test_Client_Single_Backup(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
		}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

	time.Sleep(2 * time.Second)

	{
		err := cli.Backup().Create()
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
		err := cli.Simple().Create().Element("foo", "bar")
		if err != nil {
			t.Fatal(err)
		}
	}

	time.Sleep(2 * time.Second)

	{
		err := cli.Backup().Create()
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
}
