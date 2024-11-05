package delete

type Fake struct {
	FakeClean func() error
	FakeIndex func() error
	FakeLimit func() error
	FakeScore func() error
	FakeValue func() error
}

func (f *Fake) Clean(key string) error {
	if f.FakeClean != nil {
		return f.FakeClean()
	}

	return nil
}

func (f *Fake) Index(key string, val ...string) error {
	if f.FakeIndex != nil {
		return f.FakeIndex()
	}

	return nil
}

func (f *Fake) Limit(key string, lim int) error {
	if f.FakeLimit != nil {
		return f.FakeLimit()
	}

	return nil
}

func (f *Fake) Score(key string, sco float64, end ...float64) error {
	if f.FakeScore != nil {
		return f.FakeScore()
	}

	return nil
}

func (f *Fake) Value(key string, val ...string) error {
	if f.FakeValue != nil {
		return f.FakeValue()
	}

	return nil
}
