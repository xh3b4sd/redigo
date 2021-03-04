// +build single

package client

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/redigo"
)

func Test_Client_Single_Locker_Lifecycle(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := Config{
			Kind: KindSingle,
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

	erc := make(chan error, 1)

	go func() {
		defer close(erc)

		var s string
		var w sync.WaitGroup

		w.Add(2)

		{
			err = cli.Sorted().Create().Element("key", "val", 1)
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}
		}

		go func() {
			defer w.Done()

			err = cli.Locker().Acquire()
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}

			str, err := cli.Sorted().Search().Order("key", 0, 1)
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}

			l := len(strings.Split(str[0], ":"))
			if l == 1 {
				s = "a"
				_, err := cli.Sorted().Update().Value("key", str[0]+":a", 1)
				if err != nil {
					erc <- tracer.Mask(err)
					return
				}
			}

			err = cli.Locker().Release()
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}
		}()

		go func() {
			defer w.Done()

			err = cli.Locker().Acquire()
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}

			str, err := cli.Sorted().Search().Order("key", 0, 1)
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}

			l := len(strings.Split(str[0], ":"))
			if l == 1 {
				s = "b"
				_, err := cli.Sorted().Update().Value("key", str[0]+":b", 1)
				if err != nil {
					erc <- tracer.Mask(err)
					return
				}
			}

			err = cli.Locker().Release()
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}
		}()

		w.Wait()

		{
			str, err := cli.Sorted().Search().Order("key", 0, 1)
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}

			l := strings.Split(str[0], ":")
			if len(l) != 2 {
				erc <- fmt.Errorf("l must be 2")
				return
			}
			if s == "a" && l[1] != "a" {
				erc <- fmt.Errorf("locking failed")
				return
			}
			if s == "b" && l[1] != "b" {
				erc <- fmt.Errorf("locking failed")
				return
			}
		}

		{
			err = cli.Sorted().Delete().Score("key", 1)
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}
		}
	}()

	{
		err = <-erc
		if err != nil {
			t.Fatal(err)
		}
	}
}
