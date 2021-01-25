package pubsub

import "github.com/xh3b4sd/tracer"

func (p *PubSub) Pub(key string, val string) error {
	con := p.pool.Get()
	defer con.Close()

	_, err := con.Do("PUBLISH", key, val)
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
