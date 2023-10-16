package fake

type SortedFloats struct {
	FakeScore func() (float64, error)
}

func (f *SortedFloats) Score(key string, val string, sco float64) (float64, error) {
	if f.FakeScore != nil {
		return f.FakeScore()
	}

	return 0, nil
}
