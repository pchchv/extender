package httpext

import (
	"net/url"

	"github.com/pchchv/form"
)

var (
	// DefaultFormEncoder of this package, which is configurable.
	DefaultFormEncoder FormEncoder = form.NewEncoder()
	// DefaultFormDecoder of this package, which is configurable.
	DefaultFormDecoder FormDecoder = form.NewDecoder()
)

// FormEncoder is the type used for encoding form data.
type FormEncoder interface {
	Encode(interface{}) (url.Values, error)
}

// FormDecoder is the type used for decoding a form for use.
type FormDecoder interface {
	Decode(interface{}, url.Values) error
}
