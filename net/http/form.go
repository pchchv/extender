package httpext

import (
	"net/url"

	"github.com/pchchv/form"
)

// DefaultFormEncoder of this package, which is configurable.
var DefaultFormEncoder FormEncoder = form.NewEncoder()

// FormEncoder is the type used for encoding form data.
type FormEncoder interface {
	Encode(interface{}) (url.Values, error)
}
