package asciiext

// IsUpper returns true if the byte is an ASCII uppercase letter.
func IsUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

// IsLower returns true if the byte is an ASCII lowercase letter.
func IsLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}
