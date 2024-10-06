package walker

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/redigo/walker/search"
)

type Config struct {
	Count  int
	Pool   *redis.Pool
	Prefix string
}

type Walker struct {
	search *search.Redis
}

func New(config Config) *Walker {
	var sea *search.Redis
	{
		sea = search.New(search.Config{
			Cou: config.Count,
			Poo: config.Pool,
			Pre: config.Prefix,
		})
	}

	return &Walker{
		search: sea,
	}
}

func (w *Walker) Search() Search {
	return w.search
}
