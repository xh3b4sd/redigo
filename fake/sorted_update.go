package fake

type SortedUpdate struct {
	FakeValue func() (bool, error)
}

func (u *SortedUpdate) Value(key string, new string, sco float64, ind ...string) (bool, error) {
	if u.FakeValue != nil {
		return u.FakeValue()
	}

	return false, nil
}
