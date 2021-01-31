package walker

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

func (w *Walker) Simple(pat string, don <-chan struct{}, res chan<- string) error {
	con := w.pool.Get()
	defer con.Close()

	var cur int64

	// Start to scan the set until the cursor is 0 again. Note that we check for
	// the closer twice. At first we prevent scans in case the closer was
	// triggered directly, and second before each channel send. That way ending
	// the walk immediately is guaranteed.
	for {
		select {
		case <-don:
			return nil
		default:
		}

		val, err := redis.Values(con.Do("SCAN", cur, "MATCH", pat, "COUNT", w.count))
		if err != nil {
			return tracer.Mask(err)
		}

		num, str, err := decode(val)
		if err != nil {
			return tracer.Mask(err)
		}
		cur = num

		for _, s := range str {
			select {
			case <-don:
				return nil
			default:
			}

			res <- s
		}

		if cur == 0 {
			break
		}
	}

	return nil
}
