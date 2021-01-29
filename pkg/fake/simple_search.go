package fake

type SimpleSearch struct {
	FakeValue func() (string, error)
}

func (s *SimpleSearch) Value(key string) (string, error) {
	if s.FakeValue != nil {
		return s.FakeValue()
	}

	return "", nil
}
