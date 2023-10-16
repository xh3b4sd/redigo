package fake

type Walker struct {
	FakeSimple func(pat string, don <-chan struct{}, key chan<- string) error
}

func (f *Walker) Simple(pat string, don <-chan struct{}, key chan<- string) error {
	if f.FakeSimple != nil {
		return f.FakeSimple(pat, don, key)
	}

	return nil
}
