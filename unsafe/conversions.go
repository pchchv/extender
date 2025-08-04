package unsafeext

import "unsafe"

// BytesToString converts an array of bytes into a string without allocating.
// The byte slice passed to this function is not to be used after this call as it's unsafe.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
