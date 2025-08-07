package httpext

import (
	"context"
	"net/http"
	"strconv"

	resultext "github.com/pchchv/extender/values/result"
)

// BuildRequestFn is a function used to rebuild an HTTP request for use in retryable code.
type BuildRequestFn func(ctx context.Context) resultext.Result[*http.Request, error]

// IsRetryableStatusCodeFn is a function used to determine if
// the provided status code is considered retryable.
type IsRetryableStatusCodeFn func(ctx context.Context, code int) bool

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
