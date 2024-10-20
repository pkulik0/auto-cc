package errs

import "errors"

var (
	// NotFound is returned when a resource is not found.
	NotFound = errors.New("autocc: not found")
	// InvalidInput is returned when an input is invalid.
	InvalidInput = errors.New("autocc: invalid input")
	// SourceClosedCaptionsNotFound is returned when the source closed captions are not found.
	SourceClosedCaptionsNotFound = errors.New("autocc: source closed captions not found")
)
