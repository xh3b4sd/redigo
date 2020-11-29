package fake

type SortedDelete struct {
	FakeElement func() error
}

func (d *SortedDelete) Element(key string, val string, ind ...string) error {
	if d.FakeElement != nil {
		return d.FakeElement()
	}

	return nil
}
