package backup

type Interface interface {
	// Create initiates a snapshot and blocks until it is complete.
	//
	//     https://redis.io/commands/bgsave
	//
	Create() error
}
