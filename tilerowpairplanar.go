package tileconv

// TileRowPairPlanar is a Codec that encodes each tile as a planar
// image, with the planes stored in pairs (as if it was a sequence of
// 2bpp tiles), and each plane pair stored in a row-planar manner.
//
// Thus: r0p0, r0p1, r1p0, r1p1, ..., r7p1, r0p2, r0p3, r1p2, ..., r7pN.
//
// This codec only really supports even bit depths; odd bit depths will
// be rounded up to the next even bit depth, except that the extra bit
// plane will be forced to zero.
type TileRowPairPlanar struct {
	BitDepth BitDepth
}

var _ Codec = TileRowPairPlanar{}

// Size implements Codec, returning the size of a tile.
func (c TileRowPairPlanar) Size() int {
	return ((c.BitDepth + 1) &^ 1).BytesPerTile()
}

// Encode implements Codec, encoding a tile image into bytes.
func (c TileRowPairPlanar) Encode(s SourceImage, x, y int, d []byte) {
	planes := c.BitDepth.Planes()
	mask := c.BitDepth.ColorMask()
	for iy := 0; iy < 8; iy++ {
		for ix := 0; ix < 8; ix++ {
			color := s.ColorIndexAt(x+ix, y+iy)
			color &= mask
			for p := 0; p < planes; p += 2 {
				i := iy*2 + p*BytesPerPlane + 0
				d[i] = (d[i] << 1) | (color & 1)
				color >>= 1

				i = iy*2 + p*BytesPerPlane + 1
				d[i] = (d[i] << 1) | (color & 1)
				color >>= 1
			}
		}
	}
}

// Decode implements Codec, decoding bytes into an image.
func (c TileRowPairPlanar) Decode(src []byte, dst DestImage, x, y int) {
	planes := c.BitDepth.Planes()
	mask := c.BitDepth.ColorMask()
	for iy := 0; iy < 8; iy++ {
		row := [8]uint8{}
		for p := 0; p < planes; p += 2 {
			d := src[iy*2+p*BytesPerPlane+0]
			for ix := 8 - 1; ix >= 0; ix-- {
				row[ix] |= (d & 1) << (p + 0)
				d >>= 1
			}

			d = src[iy*2+p*BytesPerPlane+1]
			for ix := 8 - 1; ix >= 0; ix-- {
				row[ix] |= (d & 1) << (p + 1)
				d >>= 1
			}
		}
		for ix := 0; ix < 8; ix++ {
			dst.SetColorIndex(x+ix, y+iy, row[ix]&mask)
		}
	}
}
