package delete

type Fake struct {
	FakeMulti func() (int64, error)
}

func (d *Fake) Multi(key ...string) (int64, error) {
	if d.FakeMulti != nil {
		return d.FakeMulti()
	}

	return 0, nil
}
