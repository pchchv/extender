package httpext

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	asciiext "github.com/pchchv/extender/ascii"
	. "github.com/pchchv/extender/values/option"
)

const (
	QueryParams QueryParamsOption = iota
	NoQueryParams
)

// QueryParamsOption represents the options for including query parameters during Decode helper functions.
type QueryParamsOption uint8

// HasRetryAfter parses the Retry-After header and returns the duration if possible.
func HasRetryAfter(headers http.Header) Option[time.Duration] {
	if ra := headers.Get(RetryAfter); ra != "" {
		if asciiext.IsDigit(ra[0]) {
			if n, err := strconv.ParseInt(ra, 10, 64); err == nil {
				return Some(time.Duration(n) * time.Second)
			}
		} else {
			// not a number so must be a date in the future
			if t, err := http.ParseTime(ra); err == nil {
				return Some(time.Until(t))
			}
		}
	}
	return None[time.Duration]()
}

// DecodeQueryParams takes the URL Query params flag.
func DecodeQueryParams(r *http.Request, v interface{}) (err error) {
	return decodeQueryParams(r.URL.Query(), v)
}

func decodeQueryParams(values url.Values, v interface{}) (err error) {
	err = DefaultFormDecoder.Decode(v, values)
	return
}
