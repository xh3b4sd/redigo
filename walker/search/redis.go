package search

import (
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Cou int
	Poo *redis.Pool
	Pre string
}

type Redis struct {
	cou int
	poo *redis.Pool
	pre string
}

func New(c Config) *Redis {
	if c.Cou <= 0 {
		c.Cou = 100
	}

	return &Redis{
		cou: c.Cou,
		poo: c.Poo,
		pre: c.Pre,
	}
}

func (r *Redis) Keys(pat string, don <-chan struct{}, res chan<- string) error {
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
			con := r.poo.Get()
			defer con.Close()

			defer close(erc)

			for {
				select {
				case <-don:
					return
				default:
				}

				val, err := redis.Values(con.Do("SCAN", cur, "MATCH", pat, "COUNT", r.cou))
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

func (r *Redis) Type(key string) (string, error) {
	var err error

	var con redis.Conn
	{
		con = r.poo.Get()
		defer con.Close()
	}

	var res string
	{
		res, err = redis.String(con.Do("TYPE", key))
		if err != nil {
			return "", tracer.Mask(err)
		}
	}

	if res == "" {
		return "", tracer.Mask(notFoundError)
	}

	return res, nil
}

func decode(res []interface{}) (int64, []string, error) {
	cur, err := strconv.ParseInt(string(res[0].([]uint8)), 10, 64)
	if err != nil {
		return 0, nil, tracer.Mask(err)
	}

	var str []string
	for _, v := range res[1].([]interface{}) {
		str = append(str, string(v.([]uint8)))
	}

	return cur, str, nil
}
