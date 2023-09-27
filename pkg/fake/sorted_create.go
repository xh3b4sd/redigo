package fake

type SortedCreate struct {
	FakeIndex func() error
	FakeValue func() error
}

func (e *SortedCreate) Index(key string, val string, sco float64, ind ...string) error {
	if e.FakeIndex != nil {
		return e.FakeIndex()
	}

	return nil
}

func (e *SortedCreate) Value(key string, val string, sco float64) error {
	if e.FakeValue != nil {
		return e.FakeValue()
	}

	return nil
}
