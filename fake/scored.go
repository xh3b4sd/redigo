package fake

type Scored struct {
}

func (s *Scored) Create(key string, ele string, sco float64) error {
	return nil
}

func (s *Scored) CutOff(key string, num int) error {
	return nil
}

func (s *Scored) Delete(key string, ele string) error {
	return nil
}

func (s *Scored) Search(key string, lef int, rig int) ([]string, error) {
	return nil, nil
}
