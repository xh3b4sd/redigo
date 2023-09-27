package fake

type SortedCreate struct {
	FakeIndex func() error
	FakeScore func() error
}

func (e *SortedCreate) Index(key string, val string, sco float64, ind ...string) error {
	if e.FakeIndex != nil {
		return e.FakeIndex()
	}

	return nil
}

func (e *SortedCreate) Score(key string, val string, sco float64) error {
	if e.FakeScore != nil {
		return e.FakeScore()
	}

	return nil
}
