package sorted

import (
	"github.com/gomodule/redigo/redis"
)

type Config struct {
	Pool *redis.Pool

	Prefix string
}

type Sorted struct {
	cre *create
	del *delete
	exi *exists
	flo *floats
	met *metric
	sea *search
	upd *update
}

func New(config Config) (*Sorted, error) {
	var cre *create
	{
		cre = &create{
			pool: config.Pool,

			createScoreScript: redis.NewScript(2, createScoreScript),

			prefix: config.Prefix,
		}
	}

	var del *delete
	{
		del = &delete{
			pool: config.Pool,

			deleteCleanScript: redis.NewScript(2, deleteCleanScript),
			deleteIndexScript: redis.NewScript(2, deleteIndexScript),
			deleteScoreScript: redis.NewScript(2, deleteScoreScript),

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
