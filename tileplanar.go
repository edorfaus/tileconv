package chrconv

// TilePlanar is a Codec that encodes each tile as a planar image, with
// each plane of the tile stored contiguously - such that a tile with
// bit depth N can be extended to N+1 bits by appending a zeroed plane.
type TilePlanar struct {
	BitDepth BitDepth
}

var _ Codec = TilePlanar{}

// Size implements Codec, returning the size of a tile.
func (c TilePlanar) Size() int {
	return c.BitDepth.BytesPerTile()
}

// Encode implements Codec, encoding a tile image into bytes.
func (c TilePlanar) Encode(src SourceImage, x, y int, dst []byte) {
	// We need to clear the destination, in case we don't have 8 planes.
	for i := c.BitDepth.BytesPerTile() - 1; i >= 0; i-- {
		dst[i] = 0
	}
	planes := c.BitDepth.Planes()
	for iy := 0; iy < 8; iy++ {
		for ix := 0; ix < 8; ix++ {
			color := src.ColorIndexAt(x+ix, y+iy)
			for p := 0; p < planes; p++ {
				i := iy + p*BytesPerPlane
				dst[i] = (dst[i] << 1) | (color & 1)
				color >>= 1
			}
		}
	}
}

// Decode implements Codec, decoding bytes into an image.
func (c TilePlanar) Decode(src []byte, dst DestImage, x, y int) {
	_ = src[c.BitDepth.BytesPerTile()-1]
	planes := c.BitDepth.Planes()
	for iy := 0; iy < 8; iy++ {
		row := [8]uint8{}
		for p := 0; p < planes; p++ {
			d := src[iy+p*BytesPerPlane]
			for ix := 8 - 1; ix >= 0; ix-- {
				row[ix] |= (d & 1) << p
				d >>= 1
			}
		}
		for ix := 0; ix < 8; ix++ {
			dst.SetColorIndex(x+ix, y+iy, row[ix])
		}
	}
}
