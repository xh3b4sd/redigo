package fake

type SortedDelete struct {
	FakeClean func() error
	FakeIndex func() error
	FakeLimit func() error
	FakeScore func() error
	FakeValue func() error
}

func (d *SortedDelete) Clean(key string) error {
	if d.FakeClean != nil {
		return d.FakeClean()
	}

	return nil
}

func (d *SortedDelete) Index(key string, val ...string) error {
	if d.FakeIndex != nil {
		return d.FakeIndex()
	}

	return nil
}

func (d *SortedDelete) Limit(key string, lim int) error {
	if d.FakeLimit != nil {
		return d.FakeLimit()
	}

	return nil
}

func (d *SortedDelete) Score(key string, sco float64) error {
	if d.FakeScore != nil {
		return d.FakeScore()
	}

	return nil
}

func (d *SortedDelete) Value(key string, val ...string) error {
	if d.FakeValue != nil {
		return d.FakeValue()
	}

	return nil
}
