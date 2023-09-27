package fake

type SimpleExists struct {
	FakeMulti func() (int64, error)
}

func (e *SimpleExists) Multi(key ...string) (int64, error) {
	if e.FakeMulti != nil {
		return e.FakeMulti()
	}

	return 0, nil
}
