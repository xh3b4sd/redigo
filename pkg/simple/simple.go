package simple

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/redigo/pkg/simple/create"
	"github.com/xh3b4sd/redigo/pkg/simple/delete"
)

type Config struct {
	Pool *redis.Pool

	Prefix string
}

type Simple struct {
	create *create.Redis
	delete *delete.Redis
	exists *exists
	search *search
}

func New(config Config) (*Simple, error) {
	// TODO refactor all redis methods like this new way
	var cre *create.Redis
	{
		cre = create.New(create.Config{
			Poo: config.Pool,
			Pre: config.Prefix,
		})
	}

	var del *delete.Redis
	{
		del = delete.New(delete.Config{
			Poo: config.Pool,
			Pre: config.Prefix,
		})
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
