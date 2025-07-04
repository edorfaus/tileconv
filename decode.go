package tileconv

import (
	"image"
)

// Decode all the tiles in the given byte slice into the given image,
// using the given codec to decode each tile.
//
// The tiles are placed on the destination image in row-major order from
// top to bottom, left to right; starting from the top-left corner.
//
// If the image size is not an even multiple of the tile size, then some
// of the tile graphics may be lost due to being rendered partially
// outside of the image. Similarly, if the image is not large enough to
// fit all the source tiles, then the remaining tiles will be lost.
//
// The destination image must have a palette that is large enough for
// the bit depth of the codec, otherwise this may break the image.
func Decode(src []byte, dst *image.Paletted, codec Codec) {
	b := dst.Bounds()
	sz := codec.Size()
	from, to := 0, sz
	for y := b.Min.Y; y < b.Max.Y && to <= len(src); y += 8 {
		for x := b.Min.X; x < b.Max.X && to <= len(src); x += 8 {
			codec.Decode(src[from:to], dst, x, y)
			from, to = to, to+sz
		}
	}
}
