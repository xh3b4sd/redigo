package redigo

import (
	"time"

	"github.com/xh3b4sd/breakr"
	"github.com/xh3b4sd/redigo/pkg/fake"
)

func Default() Interface {
	var err error

	var bre breakr.Interface
	{
		bre = breakr.New(breakr.Config{
			Failure: breakr.Failure{
				Budget: 30,
				Cooler: 1 * time.Second,
			},
		})
	}

	var red Interface
	{
		c := Config{
			Kind: KindSingle,
			Locker: ConfigLocker{
				Breakr: bre,
			},
		}

		red, err = New(c)
		if err != nil {
			panic(err)
		}
	}

	return red
}

func Fake() Interface {
	return fake.New()
}
