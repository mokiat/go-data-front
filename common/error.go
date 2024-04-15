package common

import "errors"

// ErrLimitsExceeded is returned when the decoder has reached the maximum
// number of allowed resources to parse.
var ErrLimitsExceeded = errors.New("safety limits exceeded")

// ErrInvalid is returned when an invalid file construct is detected.
var ErrInvalid = errors.New("invalid construct")
