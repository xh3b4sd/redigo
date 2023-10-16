package fake

import (
	"github.com/xh3b4sd/redigo/pkg/simple"
	"github.com/xh3b4sd/redigo/pkg/simple/create"
	"github.com/xh3b4sd/redigo/pkg/simple/delete"
	"github.com/xh3b4sd/redigo/pkg/simple/exists"
	"github.com/xh3b4sd/redigo/pkg/simple/search"
)

type Simple struct {
	FakeCreate func() simple.Create
	FakeDelete func() simple.Delete
	FakeExists func() simple.Exists
	FakeSearch func() simple.Search
}

func (s *Simple) Create() simple.Create {
	if s.FakeCreate != nil {
		return s.FakeCreate()
	}

	return &create.Fake{}
}

func (s *Simple) Delete() simple.Delete {
	if s.FakeDelete != nil {
		return s.FakeDelete()
	}

	return &delete.Fake{}
}

func (s *Simple) Exists() simple.Exists {
	if s.FakeExists != nil {
		return s.FakeExists()
	}

	return &exists.Fake{}
}

func (s *Simple) Search() simple.Search {
	if s.FakeSearch != nil {
		return s.FakeSearch()
	}

	return &search.Fake{}
}
