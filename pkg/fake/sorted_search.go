package fake

type SortedSearch struct {
	FakeIndex func() (string, error)
	FakeInter func() ([]string, error)
	FakeOrder func() ([]string, error)
	FakeRando func() ([]string, error)
	FakeScore func() ([]string, error)
	FakeUnion func() ([]string, error)
}

func (s *SortedSearch) Index(key string, ind string) (string, error) {
	if s.FakeIndex != nil {
		return s.FakeIndex()
	}

	return "", nil
}

func (s *SortedSearch) Inter(key ...string) ([]string, error) {
	if s.FakeInter != nil {
		return s.FakeInter()
	}

	return nil, nil
}

func (s *SortedSearch) Order(key string, lef int, rig int, sco ...bool) ([]string, error) {
	if s.FakeOrder != nil {
		return s.FakeOrder()
	}

	return nil, nil
}

func (s *SortedSearch) Rando(key string, cou ...uint) ([]string, error) {
	if s.FakeRando != nil {
		return s.FakeRando()
	}

	return nil, nil
}

func (s *SortedSearch) Score(key string, lef float64, rig float64) ([]string, error) {
	if s.FakeScore != nil {
		return s.FakeScore()
	}

	return nil, nil
}

func (s *SortedSearch) Union(key ...string) ([]string, error) {
	if s.FakeUnion != nil {
		return s.FakeUnion()
	}

	return nil, nil
}
