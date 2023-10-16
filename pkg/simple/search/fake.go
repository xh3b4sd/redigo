package search

type Fake struct {
	FakeMulti func() ([]string, error)
}

func (f *Fake) Multi(key ...string) ([]string, error) {
	if f.FakeMulti != nil {
		return f.FakeMulti()
	}

	return nil, nil
}
