package search

type Fake struct {
	FakeKeys func() error
	FakeType func() (string, error)
}

func (f *Fake) Keys(pat string, don <-chan struct{}, res chan<- string) error {
	if f.FakeKeys != nil {
		return f.FakeKeys()
	}

	return nil
}

func (f *Fake) Type(key string) (string, error) {
	if f.FakeType != nil {
		return f.FakeType()
	}

	return "", nil
}
