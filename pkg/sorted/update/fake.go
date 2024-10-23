package update

type Fake struct {
	FakeIndex func() (bool, error)
	FakeScore func() error
	FakeValue func() (bool, error)
}

func (f *Fake) Index(key string, new string, sco float64, ind ...string) (bool, error) {
	if f.FakeIndex != nil {
		return f.FakeIndex()
	}

	return false, nil
}

func (f *Fake) Score(key string, val string, sco float64) error {
	if f.FakeScore != nil {
		return f.FakeScore()
	}

	return nil
}

func (f *Fake) Value(key string, new string, sco float64) (bool, error) {
	if f.FakeValue != nil {
		return f.FakeValue()
	}

	return false, nil
}
