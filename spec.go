package redigo

type Interface interface {
	Ping() error
	Scored() Scored
	Shutdown()
	Simple() Simple
}

type Scored interface {
	Create(key string, ele string, sco float64) error
	CutOff(key string, num int) error
	Delete(key string, ele string) error
	// Search returns the list of scored elements stored under key. Note that
	// num may be -1 in order to list all elements.
	Search(key string, lef int, rig int) ([]string, error)
}

type Simple interface {
	Create(key, ele string) error
	Delete(key string) error
	Exists(key string) (bool, error)
	Search(key string) (string, error)
}
