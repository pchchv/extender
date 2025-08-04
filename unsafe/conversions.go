package unsafeext

import "unsafe"

// BytesToString converts an array of bytes into a string without allocating.
// The byte slice passed to this function is not to be used after this call as it's unsafe.
func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// StringToBytes converts an existing string into an []byte without allocating.
// The string passed to this functions is not to be used again after this call as it's unsafe.
func StringToBytes(s string) []byte {
	stringPtr := unsafe.StringData(s)
	return unsafe.Slice((*byte)(stringPtr), len(s))
}
