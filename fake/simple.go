package fake

import "github.com/xh3b4sd/redigo"

type Simple struct {
	FakeCreate func() redigo.SimpleCreate
	FakeDelete func() redigo.SimpleDelete
	FakeExists func() redigo.SimpleExists
	FakeSearch func() redigo.SimpleSearch
}

func (s *Simple) Create() redigo.SimpleCreate {
	if s.FakeCreate != nil {
		return s.FakeCreate()
	}

	return &SimpleCreate{}
}

func (s *Simple) Delete() redigo.SimpleDelete {
	if s.FakeDelete != nil {
		return s.FakeDelete()
	}

	return &SimpleDelete{}
}

func (s *Simple) Exists() redigo.SimpleExists {
	if s.FakeExists != nil {
		return s.FakeExists()
	}

	return &SimpleExists{}
}

func (s *Simple) Search() redigo.SimpleSearch {
	if s.FakeSearch != nil {
		return s.FakeSearch()
	}

	return &SimpleSearch{}
}
