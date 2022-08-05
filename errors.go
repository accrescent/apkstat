package apk

import "errors"

var (
	// BadIndex is returned when a bad index is encountered in parsing.
	ErrBadIndex = errors.New("bad index")
	// MalformedHeader is returned when an incorrect or invalid chunk header is encountered.
	ErrMalformedHeader = errors.New("malformed header")
	// InvalidChunkType is returned when an invalid chunk type is encountered.
	ErrInvalidChunkType = errors.New("encountered invalid chunk type")
	// XMLResourceNotFound is returned when a requested XML resource is not found.
	ErrResourceNotFound = errors.New("XML resource not found")
)
