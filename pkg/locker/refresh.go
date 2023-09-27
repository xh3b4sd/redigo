package locker

import (
	"github.com/xh3b4sd/budget/v3"
	"github.com/xh3b4sd/tracer"
)

func (l *Locker) Refresh() error {
	act := func() error {
		sta, err := l.mutex.Extend()
		if err != nil {
			return tracer.Mask(err)
		}

		if !sta {
			return tracer.Mask(budget.Cancel)
		}

		return nil
	}

	err := l.bud.Execute(act)
	if budget.IsCancel(err) {
		return tracer.Mask(statusError)
	} else if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
