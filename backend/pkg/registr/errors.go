package registr

import "errors"

var (
	ErrInvalidRole  = errors.New("invalid role: must be 'admin' or 'user'")
	ErrInvalidToken = errors.New("invalid token")
)
