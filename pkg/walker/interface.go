package walker

type Interface interface {
	Simple(pat string, don <-chan struct{}, res chan<- string) error
}
