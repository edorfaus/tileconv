package tileconv_test

import (
	"image"
	"image/color"
	"math/rand"
)

// testImageData is used to get test data for the codecs, that actually
// can vary a bit on each plane instead of being a regular pattern.
// (Though the regular counting pattern is also available.)
//
// This approach doesn't really feel right, as it basically consists of
// testing against another implementation made by the same person, and
// seems a bit overly complicated. However, it seems the alternative is
// either to hide it behind pregeneration, which is no better, or to do
// a prohibitive amount of manual work to set up similar test data,
// while testing fewer cases.
//
// There's probably a better way that I just haven't thought of, but,
// not having found that approach, this is the best I could do for now.
type testImageData struct {
	Pix    []uint8
	Planar [][]uint16
}

func newTestImageData(randomized bool) *testImageData {
	pix := make([]uint8, 256)
	for i := 0; i < 256; i++ {
		pix[i] = uint8(i)
	}
	if randomized {
		rand.New(rand.NewSource(0)).Shuffle(len(pix), func(i, j int) {
			pix[i], pix[j] = pix[j], pix[i]
		})
	}

	planar := make([][]uint16, 8)
	for p := 0; p < 8; p++ {
		planar[p] = make([]uint16, 16)
		planeBit := uint8(1) << p
		for y := 0; y < 16; y++ {
			v := uint16(0)
			valueBit := uint16(1) << 15
			for x := 0; x < 16; x++ {
				if pix[x+y*16]&planeBit != 0 {
					v |= valueBit
				}
				valueBit >>= 1
			}
			planar[p][y] = v
		}
	}

	return &testImageData{
		Pix:    pix,
		Planar: planar,
	}
}

func (d *testImageData) RowByte(x, y, p int) byte {
	// The way the plane data is stored, x = 4 requires no shifting,
	// x = 3 requires >> 1, and so on, until x = -4 requires >> 8.
	// a = x+4 gives us a = 0-8, with a=0 -> >> 8 and a=8 -> >> 0.
	xshift := 8 - (x + 4)
	return byte(d.Planar[p][y+4] >> xshift)
}

func (d *testImageData) Palette() color.Palette {
	pal := make(color.Palette, 256)
	for i := 0; i < 256; i++ {
		v := uint8(i)
		pal[i] = color.RGBA{R: v, G: 255 - v, B: 64, A: 255}
	}
	return pal
}

// FullImage returns a new image with content based on the test data.
func (d *testImageData) FullImage() *image.Paletted {
	img := image.NewPaletted(image.Rect(-4, -4, 8+4, 8+4), d.Palette())
	copy(img.Pix, d.Pix)
	return img
}

// BaseImage returns a new image that is intended as a base for Decode
// tests, with every pixel being different than the full image.
func (d *testImageData) BaseImage() *image.Paletted {
	img := image.NewPaletted(image.Rect(-4, -4, 8+4, 8+4), d.Palette())
	for i := 0; i < len(img.Pix); i++ {
		img.Pix[i] = ^d.Pix[i]
	}
	return img
}

func (d *testImageData) TilePlanar(x, y, planes int) []byte {
	data := make([]byte, 0, planes*8)
	for p := 0; p < planes; p++ {
		for i := 0; i < 8; i++ {
			data = append(data, d.RowByte(x, y+i, p))
		}
	}
	return data
}

func (d *testImageData) RowPlanar(x, y, planes int) []byte {
	data := make([]byte, 0, planes*8)
	for i := 0; i < 8; i++ {
		for p := 0; p < planes; p++ {
			data = append(data, d.RowByte(x, y+i, p))
		}
	}
	return data
}

func (d *testImageData) Packed(x, y, depth int) []byte {
	// This is a somewhat awkward and slow implementation, but it's
	// important that it's different than the target implementation.
	data := make([]byte, 0, depth*8)
	for yi := 0; yi < 8; yi++ {
		rowb := [8]byte{}
		for plane := depth - 1; plane >= 0; plane-- {
			b := d.RowByte(x, y+yi, plane)
			for i := 0; i < 8; i++ {
				rowb[i] = (rowb[i] << 1) | (b >> 7)
				b <<= 1
			}
		}

		row := uint64(0)
		for i := 0; i < 8; i++ {
			row = (row << depth) | uint64(rowb[i])
		}

		for i := 1; i <= depth; i++ {
			data = append(data, byte(row>>((depth-i)*8)))
		}
	}
	return data
}
