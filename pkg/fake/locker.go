package fake

type Locker struct {
	FakeAcquire func() error
	FakeRefresh func() error
	FakeRelease func() error
}

func (l *Locker) Acquire() error {
	if l.FakeAcquire != nil {
		return l.FakeAcquire()
	}

	return nil
}

func (l *Locker) Refresh() error {
	if l.FakeRefresh != nil {
		return l.FakeRefresh()
	}

	return nil
}

func (l *Locker) Release() error {
	if l.FakeRelease != nil {
		return l.FakeRelease()
	}

	return nil
}
