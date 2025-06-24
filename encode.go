package tileconv

import (
	"image"
	"io"
)

// Encode all the tiles in the given image into the given writer, using
// the given codec to encode each tile.
//
// The tiles are written in row-major order from top to bottom, left to
// right.
//
// If the image size is not an even multiple of the tile size, then the
// size is rounded up by adding to the right and bottom coordinates.
// This will make Encode ask the image for pixels outside of its bounds,
// which typically returns a default color index (usually 0).
func Encode(src image.PalettedImage, dst io.Writer, c Codec) error {
	buf := make([]byte, c.Size())
	b := src.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y += 8 {
		for x := b.Min.X; x < b.Max.X; x += 8 {
			c.Encode(src, x, y, buf)
			_, err := dst.Write(buf)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
