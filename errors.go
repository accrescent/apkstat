package apk

import "errors"

var (
	// BadIndex is returned when a bad index is encountered in parsing.
	BadIndex = errors.New("bad index")
	// MalformedHeader is returned when an incorrect or invalid chunk header is encountered.
	MalformedHeader = errors.New("malformed header")
	// InvalidChunkType is returned when an invalid chunk type is encountered.
	InvalidChunkType = errors.New("encountered invalid chunk type")
	// XMLResourceNotFound is returned when a requested XML resource is not found.
	XMLResourceNotFound = errors.New("XML resource not found")
)
