package fake

import (
	"github.com/xh3b4sd/redigo/pkg/simple"
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

	return &SimpleCreate{}
}

func (s *Simple) Delete() simple.Delete {
	if s.FakeDelete != nil {
		return s.FakeDelete()
	}

	return &SimpleDelete{}
}

func (s *Simple) Exists() simple.Exists {
	if s.FakeExists != nil {
		return s.FakeExists()
	}

	return &SimpleExists{}
}

func (s *Simple) Search() simple.Search {
	if s.FakeSearch != nil {
		return s.FakeSearch()
	}

	return &SimpleSearch{}
}
