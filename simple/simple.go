package simple

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/redigo/simple/create"
	"github.com/xh3b4sd/redigo/simple/delete"
	"github.com/xh3b4sd/redigo/simple/exists"
	"github.com/xh3b4sd/redigo/simple/search"
)

type Config struct {
	Pool *redis.Pool

	Prefix string
}

type Simple struct {
	create *create.Redis
	delete *delete.Redis
	exists *exists.Redis
	search *search.Redis
}

func New(config Config) (*Simple, error) {
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

	var exi *exists.Redis
	{
		exi = exists.New(exists.Config{
			Poo: config.Pool,
			Pre: config.Prefix,
		})
	}

	var sea *search.Redis
	{
		sea = search.New(search.Config{
			Poo: config.Pool,
			Pre: config.Prefix,
		})
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
