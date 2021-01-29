package fake

type Locker struct {
	FakeAcquire func() error
	FakeRelease func() error
}

func (l *Locker) Acquire() error {
	if l.FakeAcquire != nil {
		return l.FakeAcquire()
	}

	return nil
}

func (l *Locker) Release() error {
	if l.FakeRelease != nil {
		return l.FakeRelease()
	}

	return nil
}
