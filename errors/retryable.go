package errorsext

import "errors"

// IsRetryable returns true if the provided error is considered retryable by
// testing if it complies with an interface implementing `Retryable() bool` or
// `IsRetryable bool` and calling the function.
func IsRetryable(err error) bool {
	var t interface {
		IsRetryable() bool
	}
	if errors.As(err, &t) && t.IsRetryable() {
		return true
	}

	var t2 interface {
		Retryable() bool
	}
	return errors.As(err, &t2) && t2.Retryable()
}
