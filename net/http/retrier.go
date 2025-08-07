package httpext

import "net/http"

// ErrStatusCode can be used to treat/indicate a status code as an error and ability to indicate if it is retryable.
type ErrStatusCode struct {
	StatusCode            int         // the HTTP response status code that was encountered
	IsRetryableStatusCode bool        // indicates if the status code is considered retryable
	Headers               http.Header // contains the headers from the HTTP response
	Body                  []byte      // the optional body of the HTTP response
}
