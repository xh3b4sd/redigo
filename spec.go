package redigo

type Interface interface {
	Check() error
	Close() error
	Purge() error

	Locker() Locker
	PubSub() PubSub
	Sorted() Sorted
	Simple() Simple
}

type Locker interface {
	Acquire() error
	Release() error
}

type PubSub interface {
	Pub(key string, val string) error
	Sub(key string) (<-chan string, error)
}

type Simple interface {
	Create() SimpleCreate
	Delete() SimpleDelete
	Exists() SimpleExists
	Search() SimpleSearch
}

type Sorted interface {
	Create() SortedCreate
	Delete() SortedDelete
	Exists() SortedExists
	Search() SortedSearch
	Update() SortedUpdate
}

type SimpleCreate interface {
	Element(key, ele string) error
}

type SimpleDelete interface {
	Element(key string) error
}

type SimpleExists interface {
	// Element verifies if the given key does even exist. If the key exists and
	// true is returned, it means that there is a value of any datatype
	// associated with said key.
	Element(key string) (bool, error)
}

type SimpleSearch interface {
	Value(key string) (string, error)
}

type SortedCreate interface {
	Element(key string, val string, sco float64, ind ...string) error
}

type SortedDelete interface {
	// Score deletes the element identified by score within the specified sorted
	// set. Note that indices associated with the underlying element are purged
	// automatically as well.
	Score(key string, sco float64) error
	// Value deletes the element identified by value within the specified sorted
	// set. Note that indices must be provided as with done during creating and
	// updating the element. At this point, if indices are not provided with the
	// deletion, garbage indices will be left behind. For this reason, using
	// Score should be preferred over using Value.
	Value(key string, val string, ind ...string) error
}

type SortedExists interface {
	// Score verifies if an element with the given score exists within the
	// sorted set identified by key.
	Score(key string, sco float64) (bool, error)
	// Value verifies if an element with the given value exists within the
	// sorted set identified by key.
	Value(key string, val string) (bool, error)
}

type SortedSearch interface {
	// Index returns the list of sorted set elements stored under key. The
	// provided pointers are indices of the elements within the sorted set. Note
	// that lef must be greater than zero while not being greater than rig.
	// Further rig may be -1 in order to list all elements. The returned result
	// does not include scores, but only the values of the elements.
	Index(key string, lef int, rig int) ([]string, error)
	Score(key string, lef float64, rig float64) ([]string, error)
}

type SortedUpdate interface {
	// Update modifies the element identified by sco and sets its value to new.
	// For the sorted set implementations scores are static and must never
	// change since they get treated like unique IDs.
	Value(key string, new string, sco float64, ind ...string) (bool, error)
}
