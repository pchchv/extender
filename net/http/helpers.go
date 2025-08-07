package httpext

import (
	"compress/gzip"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	asciiext "github.com/pchchv/extender/ascii"
	ioext "github.com/pchchv/extender/io"
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

// DecodeXML decodes the request body into the provided struct and
// limits the request size via an ioext.LimitReader using the maxBytes param.
//
// The Content-Type e.g. "application/xml" and http method are not checked.
//
// NOTE: when includeQueryParams=true query params will be parsed and included e. g. route /user?test=true 'test'
// is added to parsed XML and replaces any values that may have been present
func DecodeXML(r *http.Request, qp QueryParamsOption, maxMemory int64, v interface{}) (err error) {
	var values url.Values
	if qp == QueryParams {
		values = r.URL.Query()
	}

	return decodeXML(r.Header, r.Body, qp, values, maxMemory, v)
}

func decodeQueryParams(values url.Values, v interface{}) (err error) {
	err = DefaultFormDecoder.Decode(v, values)
	return
}

func decodeXML(headers http.Header, body io.Reader, qp QueryParamsOption, values url.Values, maxMemory int64, v interface{}) error {
	if encoding := headers.Get(ContentEncoding); encoding == Gzip {
		gzr, err := gzip.NewReader(body)
		if err != nil {
			return err
		}

		defer func() {
			_ = gzr.Close()
		}()
		body = gzr
	}

	err := xml.NewDecoder(ioext.LimitReader(body, maxMemory)).Decode(v)
	if qp != QueryParams || err != nil {
		return err
	}

	return decodeQueryParams(values, v)
}
