package walker

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

func (w *Walker) Simple(pat string, don <-chan struct{}, res chan<- string) error {
	var erc chan error
	{
		erc = make(chan error)
	}

	var cur int64

	// Start to scan the set until the cursor is 0 again. Note that we check for
	// the closer twice. At first we prevent scans in case the closer was
	// triggered directly, and second before each channel send. That way ending
	// the walk immediately is guaranteed.
	{
		go func() {
			con := w.pool.Get()
			defer con.Close()

			defer close(erc)

			for {
				select {
				case <-don:
					return
				default:
				}

				val, err := redis.Values(con.Do("SCAN", cur, "MATCH", pat, "COUNT", w.count))
				if err != nil {
					erc <- tracer.Mask(err)
					return
				}

				num, str, err := decode(val)
				if err != nil {
					erc <- tracer.Mask(err)
					return
				}
				cur = num

				for _, s := range str {
					select {
					case <-don:
						return
					default:
					}

					res <- s
				}

				if cur == 0 {
					break
				}
			}
		}()
	}

	{
		select {
		case <-don:
			return nil
		case err := <-erc:
			if err != nil {
				return tracer.Mask(err)
			}
		}
	}

	return nil
}
