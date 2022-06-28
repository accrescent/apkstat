package apk

import "errors"

var (
	BadIndex         = errors.New("bad index")
	MalformedHeader  = errors.New("malformed header")
	InvalidChunkType = errors.New("encountered invalid chunk type")
)
