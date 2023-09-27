package client

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var alreadyExistsError = &tracer.Error{
	Kind: "alreadyExistsError",
}

func IsAlreadyExistsError(err error) bool {
	return errors.Is(err, alreadyExistsError)
}

var invalidConfigError = &tracer.Error{
	Kind: "invalidConfigError",
}

func IsInvalidConfig(err error) bool {
	return errors.Is(err, invalidConfigError)
}
