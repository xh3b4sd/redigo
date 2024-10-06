package fake

import (
	"github.com/xh3b4sd/redigo/walker"
	"github.com/xh3b4sd/redigo/walker/search"
)

type Walker struct {
	FakeSearch func() walker.Search
}

func (s *Walker) Search() walker.Search {
	if s.FakeSearch != nil {
		return s.FakeSearch()
	}

	return &search.Fake{}
}
