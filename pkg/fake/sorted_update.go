package fake

type SortedUpdate struct {
	FakeIndex func() (bool, error)
	FakeScore func() (bool, error)
}

func (u *SortedUpdate) Index(key string, new string, sco float64, ind ...string) (bool, error) {
	if u.FakeIndex != nil {
		return u.FakeIndex()
	}

	return false, nil
}

func (u *SortedUpdate) Score(key string, new string, sco float64) (bool, error) {
	if u.FakeScore != nil {
		return u.FakeScore()
	}

	return false, nil
}
