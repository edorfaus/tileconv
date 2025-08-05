package tileconv_test

import (
	"image"
	"image/color"
	"reflect"
	"testing"

	"github.com/edorfaus/tileconv"
)

func runCodecEncodeTests(
	t *testing.T, name string, codec tileconv.Codec,
	srcPix [][]uint8, want []byte,
) {
	t.Helper()

	runTestAt := func(t *testing.T, name string, x, y, dLen, dCap int) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Helper()

			src := newTestImage(x, y, srcPix)
			goodSrc := newTestImage(x, y, srcPix)

			// Make and fill some buffers to test for overflow.
			got := make([]byte, dCap)
			expect := make([]byte, dCap)
			for i := 0; i < len(expect); i++ {
				got[i] = byte(i) ^ 0b11001011
				expect[i] = byte(i) ^ 0b11001011
			}
			copy(expect, want)

			codec.Encode(src, x, y, got[:dLen])

			verifyImage(t, "source image corrupted", x, y, src, goodSrc)

			sz := codec.Size()
			verify(t, "bad encoded data", got[:sz], expect[:sz])
			verify(t, "corrupted dest tail", got[sz:], expect[sz:])
		})
	}

	t.Run(name, func(t *testing.T) {
		t.Helper()

		if len(want) != codec.Size() {
			// Either the calling test or codec.Size() is bad; this test
			// relies on both being correct, so fail out early.
			t.Errorf(
				"codec expects %v bytes, got %v bytes as wanted data",
				codec.Size(), len(want),
			)
			return
		}

		size := codec.Size()
		bigSize := size*2 + size/2

		// Test with destination having exact length and capacity
		runTestAt(t, "at_-2,3_capDst", -2, 3, size, size)

		// Test with destination having exact length but bigger capacity
		runTestAt(t, "at_-3,-2_lenDst", -3, -2, size, bigSize)

		// Test with big destination (longer than needed)
		runTestAt(t, "at_1,2_bigDst", 1, 2, bigSize, bigSize)
	})
}

func runCodecDecodeTests(
	t *testing.T, name string, codec tileconv.Codec,
	src []byte, wantPix [][]uint8,
) {
	t.Helper()

	runTestAt := func(t *testing.T, name string, x, y int, src []byte) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Helper()

			testSrc := make([]byte, cap(src))
			copy(testSrc, src[:cap(src)])
			testSrc = testSrc[:len(src)]

			got := newTestImage(0, 0, nil)
			codec.Decode(testSrc, got, x, y)

			sz := codec.Size()
			verify(t, "corrupted source data", testSrc[:sz], src[:sz])
			verify(
				t, "corrupted source tail",
				testSrc[sz:cap(testSrc)], src[sz:cap(src)],
			)

			want := newTestImage(x, y, wantPix)
			verifyImage(t, "output image incorrect", x, y, got, want)
		})
	}

	t.Run(name, func(t *testing.T) {
		t.Helper()

		if len(src) != codec.Size() {
			// Either the calling test or codec.Size() is bad; this test
			// relies on both being correct, so fail out early.
			t.Errorf(
				"codec expects %v bytes, got %v bytes as source data",
				codec.Size(), len(src),
			)
			return
		}

		testSrc := make([]byte, len(src)*2+len(src)/2)
		for i := 0; i < len(testSrc); i++ {
			testSrc[i] = byte(i) ^ 0b11001011
		}
		copy(testSrc, src)

		// Test with source having exact length and capacity
		exactSrc := testSrc[:len(src):len(src)]
		runTestAt(t, "at_-2,3_capSrc", -2, 3, exactSrc)

		// Test with source having exact length but bigger capacity
		runTestAt(t, "at_-3,-2_lenSrc", -3, -2, testSrc[:len(src)])

		// Test with big source (longer than needed)
		runTestAt(t, "at_1,2_bigSrc", 1, 2, testSrc)
	})
}

func verifyImage(
	t *testing.T, msg string, x, y int, got, want *image.Paletted,
) {
	t.Helper()

	if reflect.DeepEqual(got, want) {
		return
	}

	t.Errorf("%s:", msg)

	verify(t, "  Stride is different", got.Stride, want.Stride)
	verify(t, "  Rect is different", got.Rect, want.Rect)
	verify(t, "  Palette is different", got.Palette, want.Palette)

	// The max number of bad pixels that will be reported
	const maxBadPix = 4

	badPixel := func(m string, c, tx, ty int, g, w uint8) {
		t.Helper()
		if c < maxBadPix {
			t.Errorf(
				"  pixel at %v,%v (%vside): want %v, got %v",
				tx, ty, m, w, g,
			)
		}
	}

	inCount, outCount := 0, 0
	area := image.Rect(x, y, x+8, y+8)
	b := want.Bounds().Intersect(got.Bounds())
	for ty := b.Min.Y; ty < b.Max.Y; ty++ {
		for tx := b.Min.X; tx < b.Max.X; tx++ {
			g := got.ColorIndexAt(tx, ty)
			w := want.ColorIndexAt(tx, ty)

			if g == w {
				continue
			}

			if (image.Point{tx, ty}).In(area) {
				badPixel("in", inCount, tx, ty, g, w)
				inCount++
			} else {
				badPixel("out", outCount, tx, ty, g, w)
				outCount++
			}
		}
	}
	if inCount > 0 || outCount > 0 {
		t.Errorf(
			"  count of bad pixels: %v in + %v out = %v",
			inCount, outCount, inCount+outCount,
		)
	}

	pixCount := 0
	for i := 0; i < len(got.Pix) && i < len(want.Pix); i++ {
		if got.Pix[i] != want.Pix[i] {
			pixCount++
		}
	}
	if pixCount > 0 {
		t.Errorf("  Pix is different in %v entries", pixCount)
	}

	verify(t, "  Pix length is different", len(got.Pix), len(want.Pix))
}

func verify(t *testing.T, what string, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s:\nwant: %#v\n got: %#v", what, want, got)
	}
}

// pixBits takes the given pixel grid, and returns a new one of the same
// size that contains only the low N bits of the original values.
func pixBits(bits int, pix [][]uint8) [][]uint8 {
	mask := uint8(1<<bits) - 1
	out := make([][]uint8, len(pix))
	for i := 0; i < len(pix); i++ {
		out[i] = make([]uint8, len(pix[i]))
		for j := 0; j < len(pix[i]); j++ {
			out[i][j] = pix[i][j] & mask
		}
	}
	return out
}

func newTestImage(x, y int, px [][]uint8) *image.Paletted {
	img := image.NewPaletted(
		image.Rect(-4, -4, 8+4, 8+4), newTestPalette(),
	)
	for i := 0; i < len(img.Pix); i++ {
		img.Pix[i] = uint8(i) ^ 0b10101010
	}
	for i := 0; i < len(px); i++ {
		row := px[i]
		for j := 0; j < len(row); j++ {
			img.SetColorIndex(x+j, y+i, row[j])
		}
	}
	return img
}

func newTestPalette() color.Palette {
	pal := make(color.Palette, 256)
	for i := 0; i < 256; i++ {
		v := uint8(i)
		pal[i] = color.NRGBA{R: v, G: 255 - v, B: 64, A: 255}
	}
	return pal
}
