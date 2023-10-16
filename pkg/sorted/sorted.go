package sorted

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/redigo/pkg/sorted/create"
)

type Config struct {
	Pool *redis.Pool

	Prefix string
}

type Sorted struct {
	cre *create.Redis
	del *delete
	exi *exists
	flo *floats
	met *metric
	sea *search
	upd *update
}

func New(config Config) (*Sorted, error) {
	// TODO refactor all redis methods the new way
	var cre *create.Redis
	{
		cre = create.New(create.Config{
			Poo: config.Pool,
			Pre: config.Prefix,
		})
	}

	var del *delete
	{
		del = &delete{
			pool: config.Pool,

			deleteCleanScript: redis.NewScript(2, deleteCleanScript),
			deleteIndexScript: redis.NewScript(2, deleteIndexScript),

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

	var flo *floats
	{
		flo = &floats{
			pool: config.Pool,

			prefix: config.Prefix,
		}
	}

	var met *metric
	{
		met = &metric{
			pool: config.Pool,

			prefix: config.Prefix,
		}
	}

	var sea *search
	{
		sea = &search{
			pool: config.Pool,

			searchIndexScript: redis.NewScript(2, searchIndexScript),

			prefix: config.Prefix,
		}
	}

	var upd *update
	{
		upd = &update{
			pool: config.Pool,

			updateIndexScript: redis.NewScript(2, updateIndexScript),
			updateScoreScript: redis.NewScript(1, updateScoreScript),

			prefix: config.Prefix,
		}
	}

	s := &Sorted{
		cre: cre,
		del: del,
		exi: exi,
		flo: flo,
		met: met,
		sea: sea,
		upd: upd,
	}

	return s, nil
}

func (s *Sorted) Create() Create {
	return s.cre
}

func (s *Sorted) Delete() Delete {
	return s.del
}

func (s *Sorted) Exists() Exists {
	return s.exi
}

func (s *Sorted) Floats() Floats {
	return s.flo
}

func (s *Sorted) Metric() Metric {
	return s.met
}

func (s *Sorted) Search() Search {
	return s.sea
}

func (s *Sorted) Update() Update {
	return s.upd
}
