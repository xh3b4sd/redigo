package fake

type PubSub struct {
	FakePub func(key string, val string) error
	FakeSub func(key string) (<-chan string, error)
}

func (p *PubSub) Pub(key string, val string) error {
	if p.FakePub != nil {
		return p.FakePub(key, val)
	}

	return nil
}

func (p *PubSub) Sub(key string) (<-chan string, error) {
	if p.FakeSub != nil {
		return p.FakeSub(key)
	}

	return nil, nil
}
