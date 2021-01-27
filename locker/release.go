package locker

import "github.com/xh3b4sd/tracer"

func (l *Locker) Release() error {
	_, err := l.mutex.Unlock()
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
