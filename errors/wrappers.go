package errors

import "errors"

// These variables are used to give us access to existing
// functions in the std lib errors package.
var (
	As = errors.As
	Is = errors.Is
)
