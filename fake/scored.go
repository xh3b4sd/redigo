package fake

type Scored struct {
}

func (s *Scored) Create(key string, element string, score float64) error {
	return nil
}

func (s *Scored) CutOff(key string, num int) error {
	return nil
}

func (s *Scored) Delete(key string, element string) error {
	return nil
}

func (s *Scored) Search(key string, num int) ([]string, error) {
	return nil, nil
}
