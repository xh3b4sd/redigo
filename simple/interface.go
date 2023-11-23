package simple

type Interface interface {
	Create() Create
	Delete() Delete
	Exists() Exists
	Search() Search
}

type Create interface {
	// Element executes the Redis SET command, meaning Element will create a new
	// key-value pair if no value exists under the given key already, and update
	// an existing value if it exists under key.
	Element(key, val string) error
}

type Delete interface {
	Multi(key ...string) (int64, error)
}

type Exists interface {
	Multi(key ...string) (int64, error)
}

type Search interface {
	Multi(key ...string) ([]string, error)
}
