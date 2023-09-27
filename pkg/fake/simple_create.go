package fake

type SimpleCreate struct {
	FakeElement func() error
}

func (c *SimpleCreate) Element(key, element string) error {
	if c.FakeElement != nil {
		return c.FakeElement()
	}

	return nil
}
