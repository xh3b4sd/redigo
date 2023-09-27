package fake

type Walker struct {
	FakeSimple func(pat string, don <-chan struct{}, key chan<- string) error
}

func (w *Walker) Simple(pat string, don <-chan struct{}, key chan<- string) error {
	if w.FakeSimple != nil {
		return w.FakeSimple(pat, don, key)
	}

	return nil
}
