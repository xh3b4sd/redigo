//go:build single

package conformance

import (
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/tracer"
)

// prgAll is a convenience function for calling FLUSHALL. The provided redigo
// interface is returned as is.
func prgAll(red redigo.Interface) redigo.Interface {
	{
		err := red.Purge()
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	return red
}
