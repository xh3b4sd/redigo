//go:build single

package conformance

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/xh3b4sd/breakr"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/pkg/locker"
	"github.com/xh3b4sd/tracer"
)

func Test_Client_Single_Locker_Lifecycle(t *testing.T) {
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

	erc := make(chan error, 1)
	don := make(chan struct{}, 1)

	go func() {
		defer close(erc)

		var s string
		var w sync.WaitGroup

		w.Add(2)

		{
			err = cli.Sorted().Create().Index("key", "val", 1)
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

			don <- struct{}{}

			str, err := cli.Sorted().Search().Order("key", 0, 0)
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}

			l := len(strings.Split(str[0], ":"))
			if l == 1 {
				s = "a"
				_, err := cli.Sorted().Update().Index("key", str[0]+":a", 1)
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

		<-don

		go func() {
			defer w.Done()

			for {
				err = cli.Locker().Acquire()
				if locker.IsAcquire(err) {
					time.Sleep(50 * time.Millisecond)
					continue
				} else if err != nil {
					erc <- tracer.Mask(err)
					return
				}

				break
			}

			str, err := cli.Sorted().Search().Order("key", 0, 0)
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}

			l := len(strings.Split(str[0], ":"))
			if l == 1 {
				s = "b"
				_, err := cli.Sorted().Update().Index("key", str[0]+":b", 1)
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
			str, err := cli.Sorted().Search().Order("key", 0, 0)
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

func Test_Client_Single_Locker_Acquire_Budget(t *testing.T) {
	var err error

	var bre breakr.Interface
	{
		bre = breakr.New(breakr.Config{
			Failure: breakr.Failure{
				Budget: 3,
				Cooler: 1 * time.Second,
			},
		})
	}

	var cli redigo.Interface
	{
		c := redigo.Config{
			Kind: redigo.KindSingle,
			Locker: redigo.ConfigLocker{
				Breakr: bre,
				Expiry: 1 * time.Second,
			},
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

	don := make(chan struct{}, 1)

	go func() {
		err = cli.Locker().Acquire()
		if err != nil {
			panic(err)
		}

		time.Sleep(500 * time.Millisecond)

		don <- struct{}{}
	}()

	<-don

	// The first Acquire call should still hold the lock on the first try, but
	// the locker is configured with a breakr implementation that retries until
	// the lock expires and then can be acquired a second time here.
	err = cli.Locker().Acquire()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Client_Single_Locker_Acquire_Error(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := redigo.Config{
			Kind: redigo.KindSingle,
			Locker: redigo.ConfigLocker{
				Breakr: breakr.Fake(),
				Expiry: 1 * time.Second,
			},
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

	// Aquiring the lock the first time.
	{
		err = cli.Locker().Acquire()
		if err != nil {
			t.Fatal(err)
		}
	}

	// The first Acquire call should still hold the lock.
	{
		err = cli.Locker().Acquire()
		if !locker.IsAcquire(err) {
			t.Fatal("expected acquireError")
		}
	}
}

func Test_Client_Single_Locker_Acquire_Expiry(t *testing.T) {
	var err error

	var cli redigo.Interface
	{
		c := redigo.Config{
			Kind: redigo.KindSingle,
			Locker: redigo.ConfigLocker{
				Expiry: 1 * time.Second,
			},
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

	don := make(chan struct{}, 1)

	go func() {
		err = cli.Locker().Acquire()
		if err != nil {
			panic(err)
		}

		time.Sleep(2 * time.Second)

		don <- struct{}{}
	}()

	<-don

	// The first Acquire call should not hold the lock anymore due to expiry.
	err = cli.Locker().Acquire()
	if err != nil {
		t.Fatal(err)
	}
}
