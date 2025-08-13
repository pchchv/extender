package mathext

import (
	"math"
	"testing"

	. "github.com/pchchv/go-assert"
)

func TestMin(t *testing.T) {
	Equal(t, true, math.IsNaN(Min(math.NaN(), 1)))
	Equal(t, true, math.IsNaN(Min(1, math.NaN())))
	Equal(t, math.Inf(-1), Min(math.Inf(0), math.Inf(-1)))
	Equal(t, math.Inf(-1), Min(math.Inf(-1), math.Inf(0)))
	Equal(t, 1.0, Min(1.333, 1.0))
	Equal(t, 1.0, Min(1.0, 1.333))
	Equal(t, 1, Min(3, 1))
	Equal(t, 1, Min(1, 3))
	Equal(t, -0, Min(0, -0))
	Equal(t, -0, Min(-0, 0))
}
