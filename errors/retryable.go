package errorsext

import (
	"errors"
	"syscall"
)

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

// IsTemporary returns true if the provided error is
// considered retryable temporary error by testing if
// it complies with an interface implementing `Temporary() bool` and
// calling the function.
func IsTemporary(err error) bool {
	var t interface{ Temporary() bool }
	return errors.As(err, &t) && t.Temporary()
}

// IsTemporaryConnection returns if the
// provided error was a low level retryable connection error.
// It also returns the type, in string form,
// for optional logging and metrics use.
func IsTemporaryConnection(err error) (retryType string, isRetryable bool) {
	if err != nil {
		if errors.Is(err, syscall.ECONNRESET) {
			return "econnreset", true
		}

		if errors.Is(err, syscall.ECONNABORTED) {
			return "econnaborted", true
		}

		if errors.Is(err, syscall.ENOTCONN) {
			return "enotconn", true
		}

		if errors.Is(err, syscall.EWOULDBLOCK) {
			return "ewouldblock", true
		}

		if errors.Is(err, syscall.EAGAIN) {
			return "eagain", true
		}

		if errors.Is(err, syscall.ETIMEDOUT) {
			return "etimedout", true
		}

		if errors.Is(err, syscall.EINTR) {
			return "eintr", true
		}

		if errors.Is(err, syscall.EPIPE) {
			return "epipe", true
		}
	}
	return "", false
}

// IsTimeout returns true if the provided error is considered
// a retryable timeout error by testing if it complies with
// an interface implementing `Timeout() bool` and calling the function.
func IsTimeout(err error) bool {
	var t interface {
		Timeout() bool
	}
	return errors.As(err, &t) && t.Timeout()
}
