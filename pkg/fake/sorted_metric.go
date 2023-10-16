package fake

type SortedMetric struct {
	FakeCount func() (int64, error)
}

func (f *SortedMetric) Count(key string) (int64, error) {
	if f.FakeCount != nil {
		return f.FakeCount()
	}

	return 0, nil
}
