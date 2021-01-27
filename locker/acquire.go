package locker

import "github.com/xh3b4sd/tracer"

func (l *Locker) Acquire() error {
	err := l.mutex.Lock()
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
