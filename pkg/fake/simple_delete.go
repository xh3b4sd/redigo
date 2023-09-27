package fake

type SimpleDelete struct {
	FakeMulti func() (int64, error)
}

func (d *SimpleDelete) Multi(key ...string) (int64, error) {
	if d.FakeMulti != nil {
		return d.FakeMulti()
	}

	return 0, nil
}
