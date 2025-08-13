package mathext

import (
	"math"
	"testing"

	. "github.com/pchchv/go-assert"
)

func TestMax(t *testing.T) {
	Equal(t, true, math.IsNaN(Max(math.NaN(), 1)))
	Equal(t, true, math.IsNaN(Max(1, math.NaN())))
	Equal(t, math.Inf(0), Max(math.Inf(0), math.Inf(-1)))
	Equal(t, math.Inf(0), Max(math.Inf(-1), math.Inf(0)))
	Equal(t, 1.333, Max(1.333, 1.0))
	Equal(t, 1.333, Max(1.0, 1.333))
	Equal(t, 3, Max(3, 1))
	Equal(t, 3, Max(1, 3))
	Equal(t, 0, Max(0, -0))
	Equal(t, 0, Max(-0, 0))
}
