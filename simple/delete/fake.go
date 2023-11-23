package delete

type Fake struct {
	FakeMulti func() (int64, error)
}

func (f *Fake) Multi(key ...string) (int64, error) {
	if f.FakeMulti != nil {
		return f.FakeMulti()
	}

	return 0, nil
}
