package fake

type SortedSearch struct {
	FakeIndex func() (string, error)
	FakeInter func() ([]string, error)
	FakeOrder func() ([]string, error)
	FakeRando func() ([]string, error)
	FakeScore func() ([]string, error)
	FakeUnion func() ([]string, error)
}

func (f *SortedSearch) Index(key string, ind string) (string, error) {
	if f.FakeIndex != nil {
		return f.FakeIndex()
	}

	return "", nil
}

func (f *SortedSearch) Inter(key ...string) ([]string, error) {
	if f.FakeInter != nil {
		return f.FakeInter()
	}

	return nil, nil
}

func (f *SortedSearch) Order(key string, lef int, rig int, sco ...bool) ([]string, error) {
	if f.FakeOrder != nil {
		return f.FakeOrder()
	}

	return nil, nil
}

func (f *SortedSearch) Rando(key string, cou ...uint) ([]string, error) {
	if f.FakeRando != nil {
		return f.FakeRando()
	}

	return nil, nil
}

func (f *SortedSearch) Score(key string, lef float64, rig float64) ([]string, error) {
	if f.FakeScore != nil {
		return f.FakeScore()
	}

	return nil, nil
}

func (f *SortedSearch) Union(key ...string) ([]string, error) {
	if f.FakeUnion != nil {
		return f.FakeUnion()
	}

	return nil, nil
}
