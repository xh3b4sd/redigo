package fake

type SimpleDelete struct {
	FakeElement func() error
}

func (d *SimpleDelete) Element(key string) error {
	if d.FakeElement != nil {
		return d.FakeElement()
	}

	return nil
}
