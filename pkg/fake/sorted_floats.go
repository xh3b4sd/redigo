package fake

type SortedFloats struct {
	FakeScore func() (float64, error)
}

func (e *SortedFloats) Score(key string, val string, sco float64) (float64, error) {
	if e.FakeScore != nil {
		return e.FakeScore()
	}

	return 0, nil
}
