package sorted

import (
	"github.com/gomodule/redigo/redis"

	"github.com/xh3b4sd/redigo"
)

type Config struct {
	Pool *redis.Pool

	Prefix string
}

type Sorted struct {
	create *Create
	delete *Delete
	exists *Exists
	search *Search
	update *Update
}

func New(config Config) (*Sorted, error) {
	var cre *Create
	{
		cre = &Create{
			pool: config.Pool,

			createElementScript: redis.NewScript(2, createElementScript),

			prefix: config.Prefix,
		}
	}

	var del *Delete
	{
		del = &Delete{
			pool: config.Pool,

			scoreScript: redis.NewScript(2, deleteScoreScript),
			valueScript: redis.NewScript(2, deleteValueScript),

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

	var upd *Update
	{
		upd = &Update{
			pool: config.Pool,

			updateValueScript: redis.NewScript(2, updateValueScript),

			prefix: config.Prefix,
		}
	}

	s := &Sorted{
		create: cre,
		delete: del,
		exists: exi,
		search: sea,
		update: upd,
	}

	return s, nil
}

func (s *Sorted) Create() redigo.SortedCreate {
	return s.create
}

func (s *Sorted) Delete() redigo.SortedDelete {
	return s.delete
}

func (s *Sorted) Exists() redigo.SortedExists {
	return s.exists
}

func (s *Sorted) Search() redigo.SortedSearch {
	return s.search
}

func (s *Sorted) Update() redigo.SortedUpdate {
	return s.update
}
