package fake

type SortedExists struct {
	FakeIndex func() (bool, error)
	FakeScore func() (bool, error)
	FakeValue func() (bool, error)
}

func (f *SortedExists) Index(key string, ind string) (bool, error) {
	if f.FakeIndex != nil {
		return f.FakeIndex()
	}

	return false, nil
}

func (f *SortedExists) Score(key string, sco float64) (bool, error) {
	if f.FakeScore != nil {
		return f.FakeScore()
	}

	return false, nil
}

func (f *SortedExists) Value(key string, val string) (bool, error) {
	if f.FakeValue != nil {
		return f.FakeValue()
	}

	return false, nil
}
