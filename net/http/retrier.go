package httpext

import (
	"context"
	"net/http"
	"strconv"
	"time"

	bytesext "github.com/pchchv/extender/bytes"
	errorsext "github.com/pchchv/extender/errors"
	resultext "github.com/pchchv/extender/values/result"
)

// BuildRequestFn is a function used to rebuild an HTTP request for use in retryable code.
type BuildRequestFn func(ctx context.Context) resultext.Result[*http.Request, error]

// IsRetryableStatusCodeFn is a function used to determine if
// the provided status code is considered retryable.
type IsRetryableStatusCodeFn func(ctx context.Context, code int) bool

// DecodeAnyFn is a function used to decode the response body into the desired type.
type DecodeAnyFn func(ctx context.Context, resp *http.Response, maxMemory bytesext.Bytes, v any) error

// ErrStatusCode can be used to treat/indicate a status code as an error and ability to indicate if it is retryable.
type ErrStatusCode struct {
	StatusCode            int         // the HTTP response status code that was encountered
	IsRetryableStatusCode bool        // indicates if the status code is considered retryable
	Headers               http.Header // contains the headers from the HTTP response
	Body                  []byte      // the optional body of the HTTP response
}

// Error returns the error message for the status code.
func (e ErrStatusCode) Error() string {
	return "status code encountered: " + strconv.Itoa(e.StatusCode)
}

// IsRetryable returns if the provided status code is considered retryable.
func (e ErrStatusCode) IsRetryable() bool {
	return e.IsRetryableStatusCode
}

// Retryer is used to retry any fallible operation.
//
// The `Retryer` is designed to be stateless and reusable.
// Configuration is also copy and so a base `Retryer` can
// be used and changed for one-off requests.
// E. g. changing max attempts resulting in a new `Retrier` for that request.
type Retryer struct {
	isRetryableFn           errorsext.IsRetryableFn2[error]
	isRetryableStatusCodeFn IsRetryableStatusCodeFn
	isEarlyReturnFn         errorsext.EarlyReturnFn[error]
	decodeFn                DecodeAnyFn
	backoffFn               errorsext.BackoffFn[error]
	client                  *http.Client
	timeout                 time.Duration
	maxBytes                bytesext.Bytes
	mode                    errorsext.MaxAttemptsMode
	maxAttempts             uint8
}

// IsRetryableFn sets the `IsRetryableFn` for the `Retryer`.
func (r Retryer) IsRetryableFn(fn errorsext.IsRetryableFn2[error]) Retryer {
	r.isRetryableFn = fn
	return r
}

// IsRetryableStatusCodeFn is called to determine if the status code is retryable.
func (r Retryer) IsRetryableStatusCodeFn(fn IsRetryableStatusCodeFn) Retryer {
	if fn == nil {
		fn = func(_ context.Context, _ int) bool { return false }
	}

	r.isRetryableStatusCodeFn = fn
	return r
}

// IsEarlyReturnFn sets the `EarlyReturnFn` for the `Retryer`.
func (r Retryer) IsEarlyReturnFn(fn errorsext.EarlyReturnFn[error]) Retryer {
	r.isEarlyReturnFn = fn
	return r
}

// DecodeFn sets the decode function for the `Retryer`.
func (r Retryer) DecodeFn(fn DecodeAnyFn) Retryer {
	if fn == nil {
		fn = func(_ context.Context, _ *http.Response, _ bytesext.Bytes, _ any) error { return nil }
	}

	r.decodeFn = fn
	return r
}

// Backoff sets the backoff function for the `Retryer`.
func (r Retryer) Backoff(fn errorsext.BackoffFn[error]) Retryer {
	r.backoffFn = fn
	return r
}

// MaxAttempts sets the maximum number of attempts for the `Retryer`.
//
// NOTE: Max attempts is optional and if not set will retry indefinitely on retryable errors.
func (r Retryer) MaxAttempts(mode errorsext.MaxAttemptsMode, maxAttempts uint8) Retryer {
	r.mode, r.maxAttempts = mode, maxAttempts
	return r
}

// MaxBytes sets the maximum memory to use when decoding the response body including:
// - upon unexpected status codes.
// - when decoding the response body.
// - when draining the response body before closing allowing connection re-use.
func (r Retryer) MaxBytes(i bytesext.Bytes) Retryer {
	r.maxBytes = i
	return r
}
