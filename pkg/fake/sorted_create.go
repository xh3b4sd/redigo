package fake

type SortedCreate struct {
	FakeScore func() error
	FakeValue func() error
}

func (e *SortedCreate) Score(key string, val string, sco float64, ind ...string) error {
	if e.FakeScore != nil {
		return e.FakeScore()
	}

	return nil
}

func (e *SortedCreate) Value(key string, val string, sco float64) error {
	if e.FakeValue != nil {
		return e.FakeValue()
	}

	return nil
}
