package locker

import (
	"errors"
	"strings"

	"github.com/xh3b4sd/tracer"
)

var acquireError = &tracer.Error{
	Kind: "acquireError",
}

func IsAcquire(err error) bool {
	return errors.Is(err, acquireError) || (err != nil && strings.Contains(err.Error(), "lock already taken"))
}

var statusError = &tracer.Error{
	Kind: "statusError",
}

func IsStatus(err error) bool {
	return errors.Is(err, statusError)
}
