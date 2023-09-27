package fake

import (
	"github.com/xh3b4sd/redigo/pkg/sorted"
)

type Sorted struct {
	FakeCreate func() sorted.Create
	FakeDelete func() sorted.Delete
	FakeExists func() sorted.Exists
	FakeFloats func() sorted.Floats
	FakeMetric func() sorted.Metric
	FakeSearch func() sorted.Search
	FakeUpdate func() sorted.Update
}

func (s *Sorted) Create() sorted.Create {
	if s.FakeCreate != nil {
		return s.FakeCreate()
	}

	return &SortedCreate{}
}

func (s *Sorted) Delete() sorted.Delete {
	if s.FakeDelete != nil {
		return s.FakeDelete()
	}

	return &SortedDelete{}
}

func (s *Sorted) Exists() sorted.Exists {
	if s.FakeExists != nil {
		return s.FakeExists()
	}

	return &SortedExists{}
}

func (s *Sorted) Floats() sorted.Floats {
	if s.FakeFloats != nil {
		return s.FakeFloats()
	}

	return &SortedFloats{}
}

func (s *Sorted) Metric() sorted.Metric {
	if s.FakeMetric != nil {
		return s.FakeMetric()
	}

	return &SortedMetric{}
}

func (s *Sorted) Search() sorted.Search {
	if s.FakeSearch != nil {
		return s.FakeSearch()
	}

	return &SortedSearch{}
}

func (s *Sorted) Update() sorted.Update {
	if s.FakeUpdate != nil {
		return s.FakeUpdate()
	}

	return &SortedUpdate{}
}
