//go:build single

package conformance

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/pkg/client"
	"github.com/xh3b4sd/tracer"
)

// Test_Client_Single_Walker_Simple_001 ensures the lifecycle of scanning keys
// works as expected.
func Test_Client_Single_Walker_Simple_001(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := client.Config{
			Count: 1,
			Kind:  client.KindSingle,
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

		err = cli.Sorted().Create().Index("key", "val", 1)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Sorted().Create().Index("ssk", "foo", 0.8)
		if err != nil {
			t.Fatal(err)
		}
	}

	don := make(chan struct{}, 1)
	erc := make(chan error, 1)
	res := make(chan string, 1)

	var str []string
	var wai sync.WaitGroup
	{
		go func() {
			wai.Add(1)
			defer wai.Done()

			defer close(don)
			defer close(erc)
			defer close(res)

			go func() {
				wai.Add(1)
				defer wai.Done()

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
				wai.Add(1)
				defer wai.Done()

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

	wai.Wait()

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

// Test_Client_Single_Walker_Simple_002 ensures that pattern matching when
// scanning keys works as expected.
func Test_Client_Single_Walker_Simple_002(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := client.Config{
			Count: 1,
			Kind:  client.KindSingle,
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

	{
		// PubSub channels do not produce scannable keys.
		_, err = cli.PubSub().Sub("cha")
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Simple().Create().Element("pre:foo", "bar")
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Sorted().Create().Index("pre:key", "val", 1)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Sorted().Create().Index("ssk", "foo", 0.8)
		if err != nil {
			t.Fatal(err)
		}
	}

	don := make(chan struct{}, 1)
	erc := make(chan error, 1)
	res := make(chan string, 1)

	var str []string
	var wai sync.WaitGroup
	{
		go func() {
			wai.Add(1)
			defer wai.Done()

			defer close(don)
			defer close(erc)
			defer close(res)

			go func() {
				wai.Add(1)
				defer wai.Done()

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
				wai.Add(1)
				defer wai.Done()

				for s := range res {
					str = append(str, s)
				}
			}()

			err = cli.Walker().Simple("pre:*", don, res)
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

	wai.Wait()

	{
		if len(str) != 2 {
			t.Fatal("2 keys must be found")
		}

		for _, s := range str {
			if s == "pre:key" {
				continue
			}
			if s == "pre:foo" {
				continue
			}

			t.Fatal("key must be known")
		}
	}
}

// Test_Client_Single_Walker_Simple_003 ensures that cancelling scanning keys
// early works as expected.
func Test_Client_Single_Walker_Simple_003(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := client.Config{
			Count: 1,
			Kind:  client.KindSingle,
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

	{
		err = cli.Simple().Create().Element("foo", "bar")
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Sorted().Create().Index("key", "val", 1)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Sorted().Create().Index("ssk", "foo", 0.8)
		if err != nil {
			t.Fatal(err)
		}
	}

	don := make(chan struct{})
	res := make(chan string)

	var str []string
	var wai sync.WaitGroup
	{
		wai.Add(1)

		go func() {
			defer wai.Done()

			for {
				select {
				case <-time.After(time.Second):
					panic("test timed out")

				case <-don:
					return
				}
			}
		}()
	}

	{
		wai.Add(1)

		go func() {
			defer wai.Done()

			for s := range res {
				str = append(str, s)
				close(don)
				break
			}
		}()
	}

	{
		wai.Add(1)

		go func() {
			defer wai.Done()

			err := cli.Walker().Simple("*", don, res)
			if err != nil {
				panic(err)
			}
		}()
	}

	wai.Wait()

	{
		if len(str) != 1 {
			t.Fatal("1 keys must be found")
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
