package fake

type SimpleSearch struct {
	FakeMulti func() ([]string, error)
}

func (s *SimpleSearch) Multi(key ...string) ([]string, error) {
	if s.FakeMulti != nil {
		return s.FakeMulti()
	}

	return nil, nil
}
