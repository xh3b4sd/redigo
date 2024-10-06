package walker

type Interface interface {
	Search() Search
}

type Search interface {
	Keys(pat string, dnc <-chan struct{}, rsc chan<- string) error
	Type(key string) (string, error)
}
