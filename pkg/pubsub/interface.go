package pubsub

type Interface interface {
	Pub(key string, val string) error
	Sub(key string) (<-chan string, error)
}
