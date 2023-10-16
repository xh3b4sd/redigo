package fake

type Backup struct {
	FakeCreate func() error
}

func (f *Backup) Create() error {
	if f.FakeCreate != nil {
		return f.FakeCreate()
	}

	return nil
}
