// +build single

package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/tracer"
)

func Test_Client_Single_Walker_Simple_Lifecycle(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := Config{
			Count: 1,
			Kind:  KindSingle,
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
		// PubSub channels do not produce scannable keys.
		_, err = cli.PubSub().Sub("cha")
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Simple().Create().Element("foo", "bar")
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Sorted().Create().Element("key", "val", 1)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Sorted().Create().Element("ssk", "foo", 0.8)
		if err != nil {
			t.Fatal(err)
		}
	}

	don := make(chan struct{}, 1)
	erc := make(chan error, 1)
	res := make(chan string, 1)

	var str []string
	{
		go func() {
			defer close(don)
			defer close(erc)
			defer close(res)

			go func() {
				for {
					select {
					case <-time.After(time.Second):
						erc <- fmt.Errorf("test timed out")
						return

					case <-don:
						return
					}
				}
			}()

			go func() {
				for s := range res {
					str = append(str, s)
				}
			}()

			err = cli.Walker().Simple("*", don, res)
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}
		}()
	}

	{
		err = <-erc
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(str) != 3 {
			t.Fatal("3 keys must be found")
		}

		for _, s := range str {
			if s == "key" {
				continue
			}
			if s == "foo" {
				continue
			}
			if s == "ssk" {
				continue
			}

			t.Fatal("key must be known")
		}
	}
}

func Test_Client_Single_Walker_Simple_Cancel(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := Config{
			Count: 1,
			Kind:  KindSingle,
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
		err = cli.Simple().Create().Element("foo", "bar")
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Sorted().Create().Element("key", "val", 1)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Sorted().Create().Element("ssk", "foo", 0.8)
		if err != nil {
			t.Fatal(err)
		}
	}

	don := make(chan struct{}, 1)
	erc := make(chan error, 1)
	res := make(chan string, 1)

	var str []string
	{
		go func() {
			defer close(erc)
			defer close(res)

			go func() {
				for {
					select {
					case <-time.After(time.Second):
						erc <- fmt.Errorf("test timed out")
						return

					case <-don:
						return
					}
				}
			}()

			go func() {
				for s := range res {
					str = append(str, s)
					close(don)
				}
			}()

			err = cli.Walker().Simple("*", don, res)
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}
		}()
	}

	{
		err = <-erc
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(str) != 1 {
			t.Fatal("1 keys must be found")
		}

		for _, s := range str {
			if s == "key" {
				continue
			}

			t.Fatal("key must be known")
		}
	}
}
