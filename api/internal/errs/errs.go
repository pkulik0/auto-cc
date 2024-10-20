package errs

import "errors"

var (
	// NotFound is returned when a resource is not found.
	NotFound = errors.New("autocc: not found")
	// InvalidInput is returned when an input is invalid.
	InvalidInput = errors.New("autocc: invalid input")
)
