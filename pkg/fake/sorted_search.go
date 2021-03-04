package fake

type SortedSearch struct {
	FakeIndex func() (string, error)
	FakeOrder func() ([]string, error)
	FakeScore func() ([]string, error)
}

func (s *SortedSearch) Index(key string, ind string) (string, error) {
	if s.FakeOrder != nil {
		return s.FakeIndex()
	}

	return "", nil
}

func (s *SortedSearch) Order(key string, lef int, rig int) ([]string, error) {
	if s.FakeOrder != nil {
		return s.FakeOrder()
	}

	return nil, nil
}

func (s *SortedSearch) Score(key string, lef float64, rig float64) ([]string, error) {
	if s.FakeScore != nil {
		return s.FakeScore()
	}

	return nil, nil
}
