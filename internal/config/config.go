package config

import (
	"github.com/jinzhu/configor"
)

type Config interface{}

// Init initializes the configuration.
func Init[C Config](c *C, prefix, file string) (*C, error) {
	if err := configor.New(&configor.Config{ENVPrefix: prefix}).Load(c, file); err != nil {
		return nil, err
	}

	return c, nil
}
