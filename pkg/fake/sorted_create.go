package fake

type SortedCreate struct {
	FakeElement func() error
}

func (e *SortedCreate) Element(key string, val string, sco float64, ind ...string) error {
	if e.FakeElement != nil {
		return e.FakeElement()
	}

	return nil
}
