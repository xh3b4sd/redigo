package fake

type SortedDelete struct {
	FakeClean func() error
	FakeScore func() error
	FakeValue func() error
}

func (d *SortedDelete) Clean(key string) error {
	if d.FakeClean != nil {
		return d.FakeClean()
	}

	return nil
}

func (d *SortedDelete) Score(key string, sco float64) error {
	if d.FakeScore != nil {
		return d.FakeScore()
	}

	return nil
}

func (d *SortedDelete) Value(key string, val string) error {
	if d.FakeValue != nil {
		return d.FakeValue()
	}

	return nil
}
