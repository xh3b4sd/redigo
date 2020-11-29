package fake

type SimpleExists struct {
	FakeElement func() (bool, error)
}

func (e *SimpleExists) Element(key string) (bool, error) {
	if e.FakeElement != nil {
		return e.FakeElement()
	}

	return false, nil
}
