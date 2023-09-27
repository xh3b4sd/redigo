package fake

type SortedExists struct {
	FakeIndex func() (bool, error)
	FakeScore func() (bool, error)
	FakeValue func() (bool, error)
}

func (e *SortedExists) Index(key string, ind string) (bool, error) {
	if e.FakeIndex != nil {
		return e.FakeIndex()
	}

	return false, nil
}

func (e *SortedExists) Score(key string, sco float64) (bool, error) {
	if e.FakeScore != nil {
		return e.FakeScore()
	}

	return false, nil
}

func (e *SortedExists) Value(key string, val string) (bool, error) {
	if e.FakeValue != nil {
		return e.FakeValue()
	}

	return false, nil
}
