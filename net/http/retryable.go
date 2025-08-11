package httpext

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	errorsext "github.com/pchchv/extender/errors"
	resultext "github.com/pchchv/extender/values/result"
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

// IsRetryableStatusCodeFnR is a function used to determine if the provided status code is considered retryable.
type IsRetryableStatusCodeFnR func(code int) bool

// BuildRequestFnR is a function used to rebuild an HTTP request for use in retryable code.
type BuildRequestFnR func(ctx context.Context) (*http.Request, error)

// ErrRetryableStatusCode can be used to indicate a
// retryable HTTP status code was encountered as an error.
type ErrRetryableStatusCode struct {
	Response *http.Response
}

func (e ErrRetryableStatusCode) Error() string {
	return fmt.Sprintf("retryable HTTP status code encountered: %d", e.Response.StatusCode)
}

// ErrUnexpectedResponse can be used to indicate an unexpected response was
// encountered as an error and provide access to the *http.Response.
type ErrUnexpectedResponse struct {
	Response *http.Response
}

// IsRetryableStatusCode returns true if the provided status code is considered retryable.
func IsRetryableStatusCode(code int) bool {
	return retryableStatusCodes[code]
}

// IsNonRetryableStatusCode returns true if the provided status code should generally not be retryable.
func IsNonRetryableStatusCode(code int) bool {
	return nonRetryableStatusCodes[code]
}

// DoRetryableResponse will execute the provided functions code and automatically retry before returning the *http.Response.
//
// Deprecated: use `httpext.Retrier` instead which corrects design issues with the current implementation.
func DoRetryableResponse(ctx context.Context, onRetryFn errorsext.OnRetryFn[error], isRetryableStatusCode IsRetryableStatusCodeFnR, client *http.Client, buildFn BuildRequestFnR) resultext.Result[*http.Response, error] {
	if client == nil {
		client = http.DefaultClient
	}

	var attempt int
	for {
		req, err := buildFn(ctx)
		if err != nil {
			return resultext.Err[*http.Response, error](err)
		}

		resp, err := client.Do(req)
		if err != nil {
			if retryReason, isRetryable := errorsext.IsRetryableHTTP(err); isRetryable {
				opt := onRetryFn(ctx, err, retryReason, attempt)
				if opt.IsSome() {
					return resultext.Err[*http.Response, error](opt.Unwrap())
				}
				attempt++
				continue
			}
			return resultext.Err[*http.Response, error](err)
		}

		if isRetryableStatusCode(resp.StatusCode) {
			opt := onRetryFn(ctx, ErrRetryableStatusCode{Response: resp}, strconv.Itoa(resp.StatusCode), attempt)
			if opt.IsSome() {
				return resultext.Err[*http.Response, error](opt.Unwrap())
			}
			attempt++
			continue
		}

		return resultext.Ok[*http.Response, error](resp)
	}
}
