package fake

type Sorted struct {
	CreateFake func() error
	DeleteFake func() error
	ExistsFake func() (bool, error)
	SearchFake func() ([]string, error)
	UpdateFake func() (bool, error)
}

func (s *Sorted) Create(key string, ele string, sco float64) error {
	if s.CreateFake != nil {
		return s.CreateFake()
	}

	return nil
}

func (s *Sorted) Delete(key string, ele string) error {
	if s.DeleteFake != nil {
		return s.DeleteFake()
	}

	return nil
}

func (s *Sorted) Exists(key string, sco float64) (bool, error) {
	if s.ExistsFake != nil {
		return s.ExistsFake()
	}

	return false, nil
}

func (s *Sorted) Search(key string, lef int, rig int) ([]string, error) {
	if s.SearchFake != nil {
		return s.SearchFake()
	}

	return nil, nil
}

func (s *Sorted) Update(key string, new string, sco float64) (bool, error) {
	if s.UpdateFake != nil {
		return s.UpdateFake()
	}

	return false, nil
}
