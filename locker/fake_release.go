package locker

func (f *Fake) Release() error {
	if f.FakeRelease != nil {
		return f.FakeRelease()
	}

	return nil
}
