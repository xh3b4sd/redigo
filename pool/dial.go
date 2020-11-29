package pool

import (
	"github.com/gomodule/redigo/redis"
)

// DialConfig represents the configuration used to create a new redis
// dialer.
type DialConfig struct {
	// Address represents the address used to connect to a redis server.
	Address string
}

// DefaultDialConfig provides a default configuration to create a new
// redis dialer by best effort.
func DefaultDialConfig() DialConfig {
	newConfig := DialConfig{
		Address: "",
	}

	return newConfig
}

// NewDial creates a new configured redis dialer.
func NewDial(config DialConfig) func() (redis.Conn, error) {
	newDial := func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", config.Address)
		if err != nil {
			return nil, err
		}

		return c, nil
	}

	return newDial
}
