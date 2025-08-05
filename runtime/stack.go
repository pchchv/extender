package runtimeext

import "runtime"

// Frame wraps a runtime.Frame to provide some helper functions while
// still allowing access to the original runtime.Frame.
type Frame struct {
	runtime.Frame
}
