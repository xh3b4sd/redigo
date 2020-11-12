package fake

type Simple struct {
}

func (s *Simple) Create(key, element string) error {
	return nil
}

func (s *Simple) Delete(key string) error {
	return nil
}

func (s *Simple) Exists(key string) (bool, error) {
	return false, nil
}

func (s *Simple) Search(key string) (string, error) {
	return "", nil
}
