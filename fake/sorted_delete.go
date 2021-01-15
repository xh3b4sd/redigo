package fake

type SortedDelete struct {
	FakeScore func() error
	FakeValue func() error
}

func (d *SortedDelete) Score(key string, sco float64, ind ...string) error {
	if d.FakeScore != nil {
		return d.FakeScore()
	}

	return nil
}

func (d *SortedDelete) Value(key string, val string, ind ...string) error {
	if d.FakeValue != nil {
		return d.FakeValue()
	}

	return nil
}
