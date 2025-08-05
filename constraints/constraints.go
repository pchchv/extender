package constraintsext

// Signed represents any signed integer.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
