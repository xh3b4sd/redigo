package sorted

import (
	"errors"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

var alreadyExistsError = &tracer.Error{
	Kind: "alreadyExistsError",
}

func IsAlreadyExistsError(err error) bool {
	return errors.Is(err, alreadyExistsError)
}

var executionFailedError = &tracer.Error{
	Kind: "executionFailedError",
}

func IsExecutionFailedError(err error) bool {
	return errors.Is(err, executionFailedError)
}

var notFoundError = &tracer.Error{
	Kind: "notFoundError",
}

// IsNotFound checks whether a redis response was empty. Therefore it checks for
// redigo.ErrNil and notFoundError.
//
//     ErrNil indicates that a reply value is nil.
//
func IsNotFound(err error) bool {
	return errors.Is(err, notFoundError) || errors.Is(err, redis.ErrNil)
}
