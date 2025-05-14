package chrconv

import (
	"image"
)

// Codec is the interface implemented by each tile image format type.
//
// Each system can thus pick the codec that corresponds to the way it
// encodes its tile/sprite graphics, while reusing the surrounding code.
type Codec interface {
	// Encode the given part of the image into the given buffer.
	//
	// The x and y args specify the top-left corner of the tile in src.
	//
	// The dst slice must be large enough - at least Size() bytes long.
	//
	// Encode is not allowed to modify src, nor any part of dst that is
	// beyond the first Size() bytes.
	Encode(src image.PalettedImage, x, y int, dst []byte)

	// Decode the source data into an image at the given coordinates.
	//
	// The x and y args specify the top-left corner of the tile in dst.
	//
	// The src slice must be at least Size() bytes long; any data after
	// that is ignored.
	//
	// Decode is not allowed to modify src, nor any part of dst except
	// the color indexes inside the target area (the 8x8 tile at x,y).
	Decode(src []byte, dst SettableImage, x, y int)

	// Size returns the size of the encoded data for this codec.
	//
	// It thus gives the minimum size of any buffers that are given to
	// the Encode or Decode methods.
	Size() int
}

// SettableImage represents an indexed image that can be modified.
//
// It is the interface for images that Codec.Decode can draw into.
type SettableImage interface {
	image.Image
	SetColorIndex(x, y int, idx uint8)
}
