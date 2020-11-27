package redigo

type Interface interface {
	Ping() error
	Shutdown()
	Sorted() Sorted
	Simple() Simple
}

type Simple interface {
	Create(key, ele string) error
	Delete(key string) error
	// Exists verifies if the given key does even exist. If the key exists and
	// true is returned, it means that there is a value of any datatype
	// associated with said key.
	Exists(key string) (bool, error)
	Search(key string) (string, error)
}

type Sorted interface {
	Create() SortedCreate
	Delete() SortedDelete
	Exists() SortedExists
	Search() SortedSearch
	Update() SortedUpdate
}

type SortedCreate interface {
	Element(key string, val string, sco float64, ind ...string) error
}

type SortedDelete interface {
	Value(key string, val string) error
}

type SortedExists interface {
	// Value verifies if an element with the given value exists within the
	// sorted set identified by key.
	Value(key string, val string) (bool, error)
	// Score verifies if an element with the given score exists within the
	// sorted set identified by key.
	Score(key string, sco float64) (bool, error)
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
	// change since they get trated like unique IDs.
	Value(key string, new string, sco float64) (bool, error)
}
