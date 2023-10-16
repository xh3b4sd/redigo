package create

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var executionFailedError = &tracer.Error{
	Kind: "executionFailedError",
}

var alreadyExistsError = &tracer.Error{
	Kind: "alreadyExistsError",
}

func IsAlreadyExistsError(err error) bool {
	return errors.Is(err, alreadyExistsError)
}
