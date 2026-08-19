package postgres

import "errors"

var (
	ErrInvalidMaxPoolSize = errors.New("max pool size must be greater than 0 and less or equal 100")
)
