package chrconv

// Packed is a Codec that encodes each tile as a packed-pixel image,
// with the bits for each pixel of the tile stored contiguously, such
// that at depth 8, each byte is one pixel.
type Packed struct {
	BitDepth BitDepth
}

var _ Codec = Packed{}

// Size implements Codec, returning the size of a tile.
func (c Packed) Size() int {
	return c.BitDepth.BytesPerTile()
}

// Encode implements Codec, encoding a tile image into bytes.
func (c Packed) Encode(src SourceImage, x, y int, dst []byte) {
	bpp, mask := c.BitDepth.Planes(), c.BitDepth.ColorMask()
	di := 0
	for iy := 0; iy < 8; iy++ {
		data := byte(0)
		bits := 0
		for ix := 0; ix < 8; ix++ {
			color := mask & src.ColorIndexAt(x+ix, y+iy)
			if bits+bpp < 8 {
				data = (data << bpp) | color
				bits += bpp
			} else {
				avail := 8 - bits
				dst[di] = (data << avail) | (color >> (bpp - avail))
				di++
				data = (data << bpp) | color
				bits = bits + bpp - 8
			}
		}
	}
}

// Decode implements Codec, decoding bytes into an image.
func (c Packed) Decode(src []byte, dst DestImage, x, y int) {
	bpp := c.BitDepth.Planes()
	is := 0
	for iy := 0; iy < 8; iy++ {
		d, bits := src[is], 8
		is++

		for ix := 0; ix < 8; ix++ {
			color := d >> (8 - bpp)
			d <<= bpp
			bits -= bpp

			if bits < 0 {
				d = src[is]
				is++
				// bits missing = -bits (because bits is negative now)
				color |= d >> (8 - (-bits))
				d <<= -bits
				bits += 8
			}

			dst.SetColorIndex(x+ix, y+iy, color)
		}
	}
}
