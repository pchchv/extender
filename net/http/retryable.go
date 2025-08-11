package httpext

import (
	"fmt"
	"net/http"
)

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
	// nonRetryableStatusCodes defines common HTTP responses that are
	// not considered never to be retryable.
	nonRetryableStatusCodes = map[int]bool{
		http.StatusBadRequest:                    true,
		http.StatusUnauthorized:                  true,
		http.StatusForbidden:                     true,
		http.StatusNotFound:                      true,
		http.StatusMethodNotAllowed:              true,
		http.StatusNotAcceptable:                 true,
		http.StatusProxyAuthRequired:             true,
		http.StatusConflict:                      true,
		http.StatusLengthRequired:                true,
		http.StatusPreconditionFailed:            true,
		http.StatusRequestEntityTooLarge:         true,
		http.StatusRequestURITooLong:             true,
		http.StatusUnsupportedMediaType:          true,
		http.StatusRequestedRangeNotSatisfiable:  true,
		http.StatusExpectationFailed:             true,
		http.StatusTeapot:                        true,
		http.StatusMisdirectedRequest:            true,
		http.StatusUnprocessableEntity:           true,
		http.StatusPreconditionRequired:          true,
		http.StatusRequestHeaderFieldsTooLarge:   true,
		http.StatusUnavailableForLegalReasons:    true,
		http.StatusNotImplemented:                true,
		http.StatusHTTPVersionNotSupported:       true,
		http.StatusLoopDetected:                  true,
		http.StatusNotExtended:                   true,
		http.StatusNetworkAuthenticationRequired: true,
	}
)

// ErrRetryableStatusCode can be used to indicate a
// retryable HTTP status code was encountered as an error.
type ErrRetryableStatusCode struct {
	Response *http.Response
}

func (e ErrRetryableStatusCode) Error() string {
	return fmt.Sprintf("retryable HTTP status code encountered: %d", e.Response.StatusCode)
}

// IsRetryableStatusCodeFnR is a function used to determine if the provided status code is considered retryable.
type IsRetryableStatusCodeFnR func(code int) bool

// IsRetryableStatusCode returns true if the provided status code is considered retryable.
func IsRetryableStatusCode(code int) bool {
	return retryableStatusCodes[code]
}

// IsNonRetryableStatusCode returns true if the provided status code should generally not be retryable.
func IsNonRetryableStatusCode(code int) bool {
	return nonRetryableStatusCodes[code]
}
