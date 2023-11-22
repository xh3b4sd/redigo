package locker

import (
	"github.com/xh3b4sd/breakr"
	"github.com/xh3b4sd/tracer"
)

func (l *Locker) Refresh() error {
	act := func() error {
		sta, err := l.mut.Extend()
		if err != nil {
			return tracer.Mask(err)
		}

		if !sta {
			return tracer.Mask(breakr.Cancel)
		}

		return nil
	}

	err := l.brk.Execute(act)
	if breakr.IsCancel(err) {
		return tracer.Mask(statusError)
	} else if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
