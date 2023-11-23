package locker

func (f *Fake) Acquire() error {
	if f.FakeAcquire != nil {
		return f.FakeAcquire()
	}

	return nil
}
