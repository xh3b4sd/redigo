package simple

import (
	"github.com/gomodule/redigo/redis"

	"github.com/xh3b4sd/redigo"
)

type Config struct {
	Pool *redis.Pool

	Prefix string
}

type Simple struct {
	create *Create
	delete *Delete
	exists *Exists
	search *Search
}

func New(config Config) (*Simple, error) {
	var cre *Create
	{
		cre = &Create{
			pool: config.Pool,

			prefix: config.Prefix,
		}
	}

	var del *Delete
	{
		del = &Delete{
			pool: config.Pool,

			prefix: config.Prefix,
		}
	}

	var exi *Exists
	{
		exi = &Exists{
			pool: config.Pool,

			prefix: config.Prefix,
		}
	}

	var sea *Search
	{
		sea = &Search{
			pool: config.Pool,

			prefix: config.Prefix,
		}
	}

	s := &Simple{
		create: cre,
		delete: del,
		exists: exi,
		search: sea,
	}

	return s, nil
}

func (s *Simple) Create() redigo.SimpleCreate {
	return s.create
}

func (s *Simple) Delete() redigo.SimpleDelete {
	return s.delete
}

func (s *Simple) Exists() redigo.SimpleExists {
	return s.exists
}

func (s *Simple) Search() redigo.SimpleSearch {
	return s.search
}
