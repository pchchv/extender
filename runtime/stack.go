package runtimeext

import (
	"runtime"
	"strings"
)

// Frame wraps a runtime.Frame to provide some helper functions while
// still allowing access to the original runtime.Frame.
type Frame struct {
	runtime.Frame
}

// Function is the runtime.Frame.Function stripped down to just the function name.
func (f Frame) Function() string {
	name := f.Frame.Function
	i := strings.LastIndexByte(name, '.')
	return name[i+1:]
}
