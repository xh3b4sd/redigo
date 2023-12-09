//go:build single

package conformance

import (
	"fmt"
	"testing"
	"time"

	"github.com/xh3b4sd/redigo"
)

func Test_Client_Single_PubSub_Lifecycle(t *testing.T) {
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

	{
		cha, err := cli.PubSub().Sub("cha")
		if err != nil {
			t.Fatal(err)
		}

		go func() {
			var cou int

			defer close(erc)

			for {
				select {
				case <-time.After(time.Second):
					erc <- fmt.Errorf("test timed out")
					return

				case val := <-cha:
					cou++

					if val != "one" && val != "two" {
						erc <- fmt.Errorf("val must be one or two")
						return
					}

					if cou == 2 {
						return
					}
				}
			}
		}()
	}

	{
		err = cli.PubSub().Pub("cha", "one")
		if err != nil {
			erc <- err
		}
		err = cli.PubSub().Pub("cha", "two")
		if err != nil {
			erc <- err
		}
	}

	{
		err = <-erc
		if err != nil {
			t.Fatal(err)
		}
	}
}
