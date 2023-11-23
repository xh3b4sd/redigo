package simple

import (
	"github.com/xh3b4sd/redigo/simple/search"
)

func IsNotFound(err error) bool {
	return search.IsNotFound(err)
}
