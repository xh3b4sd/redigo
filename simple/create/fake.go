package create

type Fake struct {
	FakeElement func() error
}

func (f *Fake) Element(key, element string) error {
	if f.FakeElement != nil {
		return f.FakeElement()
	}

	return nil
}
