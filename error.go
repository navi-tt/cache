package cache

import "errors"

var (
	ErrConfig    = errors.New("wrong config type")
	InvalidKey   = errors.New("invalid key format")
	InvalidConfig = errors.New("invalid config")
	InvalidValue = errors.New("invalid value format")
	EntryNotFound =errors.New("Entry not found")

	DelFail = errors.New("del key fail")
	ErrNotFound = errors.New("key is not found")
)