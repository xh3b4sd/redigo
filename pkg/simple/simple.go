package simple

import (
	"github.com/gomodule/redigo/redis"
)

type Config struct {
	Pool *redis.Pool

	Prefix string
}

type Simple struct {
	create *create
	delete *delete
	exists *exists
	search *search
}

func New(config Config) (*Simple, error) {
	var cre *create
	{
		cre = &create{
			pool: config.Pool,

			prefix: config.Prefix,
		}
	}

	var del *delete
	{
		del = &delete{
			pool: config.Pool,

			prefix: config.Prefix,
		}
	}

	var exi *exists
	{
		exi = &exists{
			pool: config.Pool,

			prefix: config.Prefix,
		}
	}

	var sea *search
	{
		sea = &search{
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

func (s *Simple) Create() Create {
	return s.create
}

func (s *Simple) Delete() Delete {
	return s.delete
}

func (s *Simple) Exists() Exists {
	return s.exists
}

func (s *Simple) Search() Search {
	return s.search
}
