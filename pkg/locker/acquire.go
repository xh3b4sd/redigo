package locker

import (
	"github.com/xh3b4sd/tracer"
)

func (l *Locker) Acquire() error {
	act := func() error {
		err := l.mut.Lock()
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err := l.brk.Execute(act)
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
