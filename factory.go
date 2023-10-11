package redigo

import (
	"time"

	"github.com/xh3b4sd/budget/v3"
	"github.com/xh3b4sd/budget/v3/pkg/breaker"

	"github.com/xh3b4sd/redigo/pkg/fake"
)

func Default() Interface {
	var err error

	var bre budget.Interface
	{
		c := breaker.Config{
			Failure: breaker.Failure{
				Budget: 30,
				Cooler: 1 * time.Second,
			},
		}

		bre, err = breaker.New(c)
		if err != nil {
			panic(err)
		}
	}

	var red Interface
	{
		c := Config{
			Kind: KindSingle,
			Locker: ConfigLocker{
				Budget: bre,
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
