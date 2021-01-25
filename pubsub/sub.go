package pubsub

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

func (p *PubSub) Sub(key string) (<-chan string, error) {
	con := p.pool.Get()
	tic := time.NewTicker(time.Minute)

	erc := make(chan error, 1)
	mes := make(chan string, 1)
	psc := redis.PubSubConn{Conn: con}

	err := psc.Subscribe(key)
	if err != nil {
		return nil, tracer.Mask(err)
	}

	go func() {
		for {
			switch m := psc.Receive().(type) {
			case error:
				erc <- m
				return

			case redis.Message:
				mes <- string(m.Data)

			case redis.Subscription:
				switch m.Count {
				case 0:
					// All channels are unsubscribed.
					return
				case 1:
					// All channels are subscribed.
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-tic.C:
				err := psc.Ping("")
				if err != nil {
					erc <- err
				}
			case err := <-erc:
				if err != nil {
					fmt.Sprintln(err)
				}

				con.Close()
				tic.Stop()

				close(erc)
				close(mes)

				err = psc.Unsubscribe()
				if err != nil {
					fmt.Sprintln(err)
				}

				return
			}
		}

	}()

	return mes, nil
}
