package fake

type SortedUpdate struct {
	FakeIndex func() (bool, error)
	FakeScore func() (bool, error)
}

func (f *SortedUpdate) Index(key string, new string, sco float64, ind ...string) (bool, error) {
	if f.FakeIndex != nil {
		return f.FakeIndex()
	}

	return false, nil
}

func (f *SortedUpdate) Score(key string, new string, sco float64) (bool, error) {
	if f.FakeScore != nil {
		return f.FakeScore()
	}

	return false, nil
}
