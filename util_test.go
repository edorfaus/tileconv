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
	runTestAt := func(t *testing.T, name string, x, y int) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Helper()

			src := newTestImage(x, y, srcPix)
			goodSrc := newTestImage(x, y, srcPix)

			got := make([]byte, codec.Size())
			codec.Encode(src, x, y, got)

			if !reflect.DeepEqual(src, goodSrc) {
				t.Errorf(
					"source image was corrupted:\nwant: %#v\n got: %#v",
					goodSrc, src,
				)
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf(
					"bad result:\nwant: %v\n got: %v",
					want, got,
				)
			}
		})
	}
	t.Run(name, func(t *testing.T) {
		t.Helper()
		runTestAt(t, "at_-3,-2", -3, -2)
		runTestAt(t, "at_1,2", 1, 2)
	})
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
