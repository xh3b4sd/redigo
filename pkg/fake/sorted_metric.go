package fake

type SortedMetric struct {
	FakeCount func() (int64, error)
}

func (s *SortedMetric) Count(key string) (int64, error) {
	if s.FakeCount != nil {
		return s.FakeCount()
	}

	return 0, nil
}
