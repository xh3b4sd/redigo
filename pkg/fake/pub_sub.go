package fake

type PubSub struct {
	FakePub func(key string, val string) error
	FakeSub func(key string) (<-chan string, error)
}

func (f *PubSub) Pub(key string, val string) error {
	if f.FakePub != nil {
		return f.FakePub(key, val)
	}

	return nil
}

func (f *PubSub) Sub(key string) (<-chan string, error) {
	if f.FakeSub != nil {
		return f.FakeSub(key)
	}

	return nil, nil
}
