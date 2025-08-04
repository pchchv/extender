package timeext

// Instant represents a monotonic instant in time.
// Instants are opaque types that can only be compared with one another and allows measuring of duration.
type Instant int64

// NewInstant returns a new Instant.
func NewInstant() Instant {
	return Instant(NanoTime())
}
