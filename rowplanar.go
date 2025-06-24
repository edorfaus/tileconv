package tileconv

// RowPlanar is a Codec that encodes each tile as a planar image, with
// the planes for each row stored contiguously. In other words, all the
// data for the first row is stored before the data for the second, and
// so on.
type RowPlanar struct {
	BitDepth BitDepth
}

var _ Codec = RowPlanar{}

// Size implements Codec, returning the size of a tile.
func (c RowPlanar) Size() int {
	return c.BitDepth.BytesPerTile()
}

// Encode implements Codec, encoding a tile image into bytes.
func (c RowPlanar) Encode(src SourceImage, x, y int, dst []byte) {
	planes := c.BitDepth.Planes()
	for iy := 0; iy < 8; iy++ {
		for ix := 0; ix < 8; ix++ {
			color := src.ColorIndexAt(x+ix, y+iy)
			for p := 0; p < planes; p++ {
				i := iy*planes + p
				dst[i] = (dst[i] << 1) | (color & 1)
				color >>= 1
			}
		}
	}
}

// Decode implements Codec, decoding bytes into an image.
func (c RowPlanar) Decode(src []byte, dst DestImage, x, y int) {
	planes := c.BitDepth.Planes()
	for iy := 0; iy < 8; iy++ {
		row := [8]uint8{}
		for p := 0; p < planes; p++ {
			d := src[iy*planes+p]
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
