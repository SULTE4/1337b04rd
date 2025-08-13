package appError

import "errors"

var (
	ErrPostNotAvailable = errors.New("post already not available")
)
