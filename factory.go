package redigo

import (
	"github.com/xh3b4sd/redigo/pkg/fake"
)

func Default() Interface {
	var err error

	var red Interface
	{
		c := Config{
			Kind: KindSingle,
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
