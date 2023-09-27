package locker

type Interface interface {
	// Acquire creates a distributed lock.
	Acquire() error
	// Refresh prevents a distributed lock from expiring.
	Refresh() error
	// Release deletes a distributed lock so that it can be acquired by another
	// process.
	Release() error
}
