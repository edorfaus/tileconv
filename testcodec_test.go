package chrconv_test

import (
	"testing"

	"github.com/edorfaus/chrconv"
)

// testCodec is a test implementation of a Codec, used to test functions
// that take a Codec as an argument.
type testCodec struct {
	t       *testing.T
	copies  int
	encodes int
	decodes int
}

var _ chrconv.Codec = &testCodec{}

func (c *testCodec) Size() int {
	return 8 * 8 * c.copies
}

func (c *testCodec) Encode(s chrconv.SourceImage, x, y int, d []byte) {
	c.encodes++

	sz := c.Size()
	if l := len(d); l < sz {
		c.t.Errorf("dst too short in Encode: want %v, got %v", sz, l)
	}

	i := 0
	for j := c.copies - 1; j >= 0; j-- {
		for iy := 0; iy < 8; iy++ {
			for ix := 0; ix < 8 && i < sz; ix++ {
				d[i] = s.ColorIndexAt(x+ix, y+iy) << j
				i++
			}
		}
	}
}

func (c *testCodec) Decode(src []byte, dst chrconv.DestImage, x, y int) {
	c.decodes++
	// TODO
}
