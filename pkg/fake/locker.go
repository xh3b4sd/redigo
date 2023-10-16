package fake

type Locker struct {
	FakeAcquire func() error
	FakeRefresh func() error
	FakeRelease func() error
}

func (f *Locker) Acquire() error {
	if f.FakeAcquire != nil {
		return f.FakeAcquire()
	}

	return nil
}

func (f *Locker) Refresh() error {
	if f.FakeRefresh != nil {
		return f.FakeRefresh()
	}

	return nil
}

func (f *Locker) Release() error {
	if f.FakeRelease != nil {
		return f.FakeRelease()
	}

	return nil
}
