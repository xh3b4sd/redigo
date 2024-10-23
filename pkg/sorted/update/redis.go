package update

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/redigo/index"
	"github.com/xh3b4sd/redigo/prefix"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Poo *redis.Pool
	Pre string
}

type Redis struct {
	poo *redis.Pool
	pre string
	ind *redis.Script
	val *redis.Script
}

func New(c Config) *Redis {
	return &Redis{
		poo: c.Poo,
		pre: c.Pre,
		ind: redis.NewScript(2, updateIndexScript),
		val: redis.NewScript(1, updateScoreScript),
	}
}

// Index executes a script of three key operations in order to reliably modify
// the element of a sorted set. Consider the element bef being created using the
// score 23 like shown below.
//
//	redis> ZADD k:foo 25 "old"
//	(integer) 1
//
// The first step being executed is to lookup the old element associated with
// score. This is necessary because removing an element from a sorted set in our
// case does work best knowing the element itself. That way we can add the new
// element first assuming that in worst case we have two elements instead of
// zero.
//
// The second step being executed is to add the new element associated with
// score. At this point we have two elements. The old and the new one.
//
// The third step being executed is to remove the old element associated with
// score. Now the value retrieved in the first step is being leveraged. Removing
// the old element marks the end of the executed transaction, leaving a clean
// state of an updated element behind.
//
//	redis> ZRANGE k:foo 25 25 BYSCORE
//	1) "old"
//
//	redis> ZADD k:foo 25 "new"
//	(integer) 1
//
//	redis> ZREM k:foo "old"
//	(integer) 1
func (r *Redis) Index(key string, new string, sco float64, ind ...string) (bool, error) {
	var err error

	var con redis.Conn
	{
		con = r.poo.Get()
		defer con.Close()
	}

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(r.pre, key))            // KEYS[1]
		arg = append(arg, prefix.WithKeys(r.pre, index.New(key))) // KEYS[2]
		arg = append(arg, new)                                    // ARGV[1]
		arg = append(arg, sco)                                    // ARGV[2]

		for _, s := range ind {
			arg = append(arg, s)
		}
	}

	res, err := redis.Int(r.ind.Do(con, arg...))
	if err != nil {
		return false, tracer.Mask(err)
	}

	switch res {
	case 0:
		return false, tracer.Maskf(notFoundError, "sorted set does not exist under key")
	case 1:
		return false, tracer.Maskf(notFoundError, "element does not exist in sorted set")
	case 2:
		return false, nil
	case 3:
		return true, nil
	}

	return false, tracer.Mask(executionFailedError)
}

func (r *Redis) Score(key string, val string, sco float64) error {
	var err error

	var con redis.Conn
	{
		con = r.poo.Get()
		defer con.Close()
	}

	// Use XX to only update elements that already exist. Don't add new elements.
	var res int64
	{
		res, err = redis.Int64(con.Do("ZADD", prefix.WithKeys(r.pre, key), "XX", "CH", sco, val))
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// ZADD returns the number of elements updated. Scores can be duplicated, but
	// values must be unique after all. If the value to be updated does not exist,
	// then ZADD returns 0. In that case we return the appropriate error.
	if res == 0 {
		return tracer.Maskf(notFoundError, "key %s does not hold a score for value %s", key, val)
	}

	return nil
}

func (r *Redis) Value(key string, new string, sco float64) (bool, error) {
	var err error

	var con redis.Conn
	{
		con = r.poo.Get()
		defer con.Close()
	}

	var arg []interface{}
	{
		arg = append(arg, prefix.WithKeys(r.pre, key)) // KEYS[1]
		arg = append(arg, new)                         // ARGV[1]
		arg = append(arg, sco)                         // ARGV[2]
	}

	res, err := redis.Int(r.val.Do(con, arg...))
	if err != nil {
		return false, tracer.Mask(err)
	}

	switch res {
	case 0:
		return false, tracer.Maskf(notFoundError, "element does not exist in sorted set")
	case 1:
		return false, nil
	case 2:
		return true, nil
	}

	return false, tracer.Mask(executionFailedError)
}
