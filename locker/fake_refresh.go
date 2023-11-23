package locker

func (f *Fake) Refresh() error {
	if f.FakeRefresh != nil {
		return f.FakeRefresh()
	}

	return nil
}
