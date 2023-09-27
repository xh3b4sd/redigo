package sorted

type Interface interface {
	Create() Create
	Delete() Delete
	Exists() Exists
	Floats() Floats
	Metric() Metric
	Search() Search
	Update() Update
}

type Create interface {
	// Index creates an element within the sorted set under key, tracking the
	// element using the unique score given by sco. The element's value provided
	// by val can be ensured to have unique associations, like indices using ind.
	// Scores are enforced to be unique.
	Index(key string, val string, sco float64, ind ...string) error

	// Value creates an element within the sorted set under key transparently
	// using ZADD and the NX option. Scores are not enforced to be unique, values
	// are.
	//
	//     https://redis.io/commands/zadd
	//
	// TODO rename to Score due to its score interface
	Value(key string, val string, sco float64) error
}

type Delete interface {
	// Clean removes the sorted set under key including the derived indeizes.
	Clean(key string) error

	// Index deletes the elements identified by the given values within the
	// specified sorted set. Note that indices associated with the underlying
	// elements are purged automatically as well.
	Index(key string, val ...string) error

	// Limit cuts off all older elements from the sorted set under key resulting
	// in a sorted set that contains the latest lim amount of elements. Consider
	// the following elements.
	//
	//     a b c d e f g
	//
	// Executing Limit with lim set to 3 will result in the following set.
	//
	//     e f g
	//
	Limit(key string, lim int) error

	// Score deletes the element identified by the given score within the
	// specified sorted set. Non-existing elements are ignored.
	//
	//     https://redis.io/commands/zremrangebyscore
	//
	Score(key string, sco float64) error

	// Value deletes the elements identified by the given values within the
	// specified sorted set. Non-existing elements are ignored.
	//
	//     https://redis.io/commands/zrem
	//
	Value(key string, val ...string) error
}

type Exists interface {
	// Index verifies if an element with the given index exists within the
	// sorted set identified by key.
	Index(key string, ind string) (bool, error)

	// Score verifies if an element with the given score exists within the
	// sorted set identified by key.
	Score(key string, sco float64) (bool, error)

	// Value verifies if an element with the given value exists within the
	// sorted set identified by key.
	Value(key string, val string) (bool, error)
}

type Floats interface {
	// Score increments the score of the provided value in the underlying sorted
	// set by sco, in case sco is positive. Is sco negative, the associated
	// value is decremented by the provided amount. Note that elements created
	// with Floats.Score are managed separately from The Create and Update
	// interfaces, which means that element scores are not used as unique
	// identifiers.
	//
	//     https://redis.io/commands/zincrby
	//
	Score(key string, val string, sco float64) (float64, error)
}

type Metric interface {
	// Count returns the total number of elements in the underlying sorted set
	// as provided by ZCARD.
	//
	//     https://redis.io/commands/zcard
	//
	Count(key string) (int64, error)
}

type Search interface {
	// Index returns values of stored elements as associated with their indices
	// during element creation. This enables multi key elements. Values can be
	// retreived using different keys referencing the requested element's value.
	Index(key string, ind string) (string, error)

	// Inter returns the values that exist in all the given keys. Therefore the
	// returned values represent the intersection of the given keys. Given k1 and
	// k2 hold the following values, Inter(k1, k2) were to return v4 and v5.
	//
	//     k1       v3 v4 v5 v6
	//     k2    v2    v4 v5    v7
	//
	Inter(key ...string) ([]string, error)

	// Order returns the values of the sorted set elements stored under key. The
	// provided pointers are ranks of the elements' scores within the sorted set.
	// All values udner key can be returned using lef=0 and rig=-1. Optionally a
	// single bool is allowed to be passed for returning the element scores
	// instead of their values as described by WITHSCORES.
	//
	//     https://redis.io/commands/zrange
	//
	Order(key string, lef int, rig int, sco ...bool) ([]string, error)

	// Rando returns a random value within the underlying sorted set. Optionally
	// a single uint is allowed to be passed for requesting cou random values as
	// described by ZRANDMEMBER.
	//
	//     https://redis.io/commands/zrandmember
	//
	Rando(key string, cou ...uint) ([]string, error)

	// Score returns the values associated to the range of scores defined by lef
	// and rig. Can be used to find a particular value if lef and rig are equal.
	//
	//     https://redis.io/commands/zrange
	//
	Score(key string, lef float64, rig float64) ([]string, error)
}

type Update interface {
	// Index modifies the element identified by sco and sets its value to new. The
	// current implementation requires all indices to be provided that have also
	// been used to create the indexed element in the first place. For the sorted
	// set implementation here, indices and scores are static and must never
	// change since they get treated like unique IDs. The returned bool indicates
	// whether the underlying value was updated. An error is returned if the
	// underlying element does not exist.
	Index(key string, new string, sco float64, ind ...string) (bool, error)
	// Score modifies the element identified by sco and sets its value to new. For
	// the sorted set implementation here, scores are static and must never change
	// since they get treated like unique IDs. The returned bool indicates whether
	// the underlying value was updated. An error is returned if the underlying
	// element does not exist.
	Score(key string, new string, sco float64) (bool, error)
}
