package tileconv

import (
	"errors"
)

// BitDepth represents the number of bits to use per pixel of each tile.
//
// E.g. for a 4-color palette, there are 2 bits per pixel, so depth 2.
//
// Note that the methods on this type assume tiles are 8x8 pixels.
type BitDepth uint8

const (
	BD1 BitDepth = 1 + iota
	BD2
	BD3
	BD4
	BD5
	BD6
	BD7
	BD8
)

// BytesPerPlane is the number of bytes taken up by each bit plane of a
// single tile.
//
// Since each tile is 8x8 pixels, and each plane only stores 1 bit per
// pixel, each row fits in a single byte. Thus, this is just the height.
const BytesPerPlane = 8

// Planes returns the number of bit planes needed for this bit depth.
func (d BitDepth) Planes() int {
	return int(d)
}

// Colors returns the number of colors available when using this depth.
func (d BitDepth) Colors() int {
	return 1 << int(d)
}

// ColorMask returns a bitmask containing a 1 for each bit of the color
// index that is actually used when using this bit depth.
func (d BitDepth) ColorMask() uint8 {
	return (1 << d) - 1
}

// BytesPerTile returns the number of bytes that will be necessary to
// store a single tile using this bit depth.
func (d BitDepth) BytesPerTile() int {
	return int(d) * BytesPerPlane
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (d *BitDepth) UnmarshalText(text []byte) error {
	if len(text) != 1 || text[0] < '1' || text[0] > '8' {
		return errors.New("invalid bit depth")
	}
	*d = BitDepth(text[0] - '0')
	return nil
}
