package fake

import "github.com/xh3b4sd/redigo"

type Sorted struct {
	FakeCreate func() redigo.SortedCreate
	FakeDelete func() redigo.SortedDelete
	FakeExists func() redigo.SortedExists
	FakeSearch func() redigo.SortedSearch
	FakeUpdate func() redigo.SortedUpdate
}

func (s *Sorted) Create() redigo.SortedCreate {
	if s.FakeCreate != nil {
		return s.FakeCreate()
	}

	return &SortedCreate{}
}

func (s *Sorted) Delete() redigo.SortedDelete {
	if s.FakeDelete != nil {
		return s.FakeDelete()
	}

	return &SortedDelete{}
}

func (s *Sorted) Exists() redigo.SortedExists {
	if s.FakeExists != nil {
		return s.FakeExists()
	}

	return &SortedExists{}
}

func (s *Sorted) Search() redigo.SortedSearch {
	if s.FakeSearch != nil {
		return s.FakeSearch()
	}

	return &SortedSearch{}
}

func (s *Sorted) Update() redigo.SortedUpdate {
	if s.FakeUpdate != nil {
		return s.FakeUpdate()
	}

	return &SortedUpdate{}
}
