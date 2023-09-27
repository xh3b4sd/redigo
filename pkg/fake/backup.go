package fake

type Backup struct {
	FakeCreate func() error
}

func (b *Backup) Create() error {
	if b.FakeCreate != nil {
		return b.FakeCreate()
	}

	return nil
}
