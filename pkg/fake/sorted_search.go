package fake

type SortedSearch struct {
	FakeIndex func() (string, error)
	FakeInter func() ([]string, error)
	FakeOrder func() ([]string, error)
	FakeRando func() ([]string, error)
	FakeValue func() ([]string, error)
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

func (s *SortedSearch) Value(key string, lef float64, rig float64) ([]string, error) {
	if s.FakeValue != nil {
		return s.FakeValue()
	}

	return nil, nil
}
