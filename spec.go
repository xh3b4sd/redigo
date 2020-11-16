package redigo

type Interface interface {
	Ping() error
	Scored() Scored
	Shutdown()
	Simple() Simple
}

type Scored interface {
	Create(key string, ele string, sco float64) error
	Delete(key string, ele string) error
	// Exists verifies if there does any value associated with key exists.
	Exists(key string) (bool, error)
	// Search returns the list of scored elements stored under key. Note that
	// lef must be greater than zero while not being greater than rig. Further
	// rig may be -1 in order to list all elements. The returned result does not
	// include scores, but only the names of the elements stored.
	Search(key string, lef int, rig int) ([]string, error)
	// Update modifies the element identified by sco and sets its value to new.
	Update(key string, new string, sco float64) (bool, error)
}

type Simple interface {
	Create(key, ele string) error
	Delete(key string) error
	// Exists verifies if there does any value associated with key exists.
	Exists(key string) (bool, error)
	Search(key string) (string, error)
}
