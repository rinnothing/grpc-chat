package connections

import "errors"

var (
	ErrAlreadyConnected = errors.New("already connected")
	ErrNotConnected = errors.New("not connected")
)
