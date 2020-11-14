package fake

type Simple struct {
	CreateFake func() error
	DeleteFake func() error
	ExistsFake func() (bool, error)
	SearchFake func() (string, error)
}

func (s *Simple) Create(key, element string) error {
	if s.CreateFake != nil {
		return s.CreateFake()
	}

	return nil
}

func (s *Simple) Delete(key string) error {
	if s.DeleteFake != nil {
		return s.DeleteFake()
	}

	return nil
}

func (s *Simple) Exists(key string) (bool, error) {
	if s.ExistsFake != nil {
		return s.ExistsFake()
	}

	return false, nil
}

func (s *Simple) Search(key string) (string, error) {
	if s.SearchFake != nil {
		return s.SearchFake()
	}

	return "", nil
}
