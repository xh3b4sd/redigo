package fake

type Scored struct {
	CreateFake func() error
	CutOffFake func() error
	DeleteFake func() error
	SearchFake func() ([]string, error)
}

func (s *Scored) Create(key string, ele string, sco float64) error {
	if s.CreateFake != nil {
		return s.CreateFake()
	}

	return nil
}

func (s *Scored) CutOff(key string, num int) error {
	if s.CutOffFake != nil {
		return s.CutOffFake()
	}

	return nil
}

func (s *Scored) Delete(key string, ele string) error {
	if s.DeleteFake != nil {
		return s.DeleteFake()
	}

	return nil
}

func (s *Scored) Search(key string, lef int, rig int) ([]string, error) {
	if s.SearchFake != nil {
		s.SearchFake()
	}

	return nil, nil
}
