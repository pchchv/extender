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

// File is the runtime.Frame.File stripped down to just the filename.
func (f Frame) File() string {
	name := f.Frame.File
	i := strings.LastIndexByte(name, '/')
	return name[i+1:]
}

// Line is the line of the runtime.Frame and exposed for convenience.
func (f Frame) Line() int {
	return f.Frame.Line
}
