package httpext

import "net/http"

var (
	// retryableStatusCodes defines the common HTTP response codes that are considered retryable.
	retryableStatusCodes = map[int]bool{
		http.StatusServiceUnavailable: true,
		http.StatusTooManyRequests:    true,
		http.StatusBadGateway:         true,
		http.StatusGatewayTimeout:     true,
		http.StatusRequestTimeout:     true,
		// 524 is a Cloudflare specific error which indicates it connected to the
		// origin server but did not receive response within 100 seconds and so times out.
		524: true,
	}
)

// IsRetryableStatusCode returns true if the provided status code is considered retryable.
func IsRetryableStatusCode(code int) bool {
	return retryableStatusCodes[code]
}
