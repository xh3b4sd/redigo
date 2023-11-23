package locker

type Fake struct {
	FakeAcquire func() error
	FakeRefresh func() error
	FakeRelease func() error
}
