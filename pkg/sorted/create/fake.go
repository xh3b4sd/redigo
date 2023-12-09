package create

type Fake struct {
	FakeIndex func() error
	FakeScore func() error
	FakeUnion func() (int64, error)
}

func (f *Fake) Index(key string, val string, sco float64, ind ...string) error {
	if f.FakeIndex != nil {
		return f.FakeIndex()
	}

	return nil
}

func (f *Fake) Score(key string, val string, sco float64) error {
	if f.FakeScore != nil {
		return f.FakeScore()
	}

	return nil
}

func (f *Fake) Union(dst string, key ...string) (int64, error) {
	if f.FakeUnion != nil {
		return f.FakeUnion()
	}

	return 0, nil
}
