package redigo

type Interface interface {
	Ping() error
	Scored() Scored
	Shutdown()
	Simple() Simple
}

type Scored interface {
	Create(key string, element string, score float64) error
	CutOff(key string, num int) error
	Delete(key string, element string) error
	Search(key string, num int) ([]string, error)
}

type Simple interface {
	Create(key, element string) error
	Delete(key string) error
	Exists(key string) (bool, error)
	Search(key string) (string, error)
}
