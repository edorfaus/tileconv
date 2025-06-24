package tileconv_test

import (
	"image"
	"reflect"
	"testing"

	"github.com/edorfaus/tileconv"
)

func TestPackedSize(t *testing.T) {
	check := func(bd tileconv.BitDepth, want int) {
		t.Helper()
		c := tileconv.Packed{BitDepth: bd}
		got := c.Size()
		if got != want {
			t.Errorf("depth %v size: want %v, got %v", bd, want, got)
		}
	}
	check(tileconv.BD1, 8*1)
	check(tileconv.BD2, 8*2)
	check(tileconv.BD3, 8*3)
	check(tileconv.BD4, 8*4)
	check(tileconv.BD5, 8*5)
	check(tileconv.BD6, 8*6)
	check(tileconv.BD7, 8*7)
	check(tileconv.BD8, 8*8)
}

func TestPackedEncode(t *testing.T) {
	var td *testImageData
	var goodSrc *image.Paletted
	var failCount int
	const maxFails = 4

	check := func(t *testing.T, bd tileconv.BitDepth, x, y int) {
		want := td.Packed(x, y, bd.Planes())
		src := td.FullImage()
		got := make([]byte, len(want))

		c := tileconv.Packed{BitDepth: bd}
		c.Encode(src, x, y, got)

		if !reflect.DeepEqual(src, goodSrc) {
			t.Errorf(
				"BD%v @ %v,%v: source image was corrupted",
				bd, x, y,
			)
			failCount++
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf(
				"BD%v @ %v,%v: bad result:\nwant: %v\n got: %v",
				bd, x, y, want, got,
			)
			failCount++
		}
	}

	runTestForDepth := func(t *testing.T, bd tileconv.BitDepth) {
		failCount = 0
		for y := -4; y < 4 && failCount < maxFails; y++ {
			for x := -4; x < 4 && failCount < maxFails; x++ {
				check(t, bd, x, y)
			}
		}
		if failCount >= maxFails {
			t.Logf("too many errors, skipping rest of depth %v", bd)
		}
	}

	names := []string{"Seq", "Rng"}
	for i, rng := range []bool{false, true} {
		t.Run(names[i], func(t *testing.T) {
			td = newTestImageData(rng)
			goodSrc = td.FullImage()

			if src := td.FullImage(); !reflect.DeepEqual(src, goodSrc) {
				t.Errorf("bad test: source images are not equal")
				return
			}

			for bd := tileconv.BD1; bd <= tileconv.BD8; bd++ {
				runTestForDepth(t, bd)
			}
		})
	}
}

func TestPackedDecode(t *testing.T) {
	var td *testImageData
	var fullImage, baseImage *image.Paletted
	var failCount int
	const maxBadPix = 2
	const maxFails = 2

	check := func(t *testing.T, bd tileconv.BitDepth, x, y int) {
		failed := false
		errorf := func(f string, a ...any) {
			args := append([]any{bd, x, y}, a...)
			t.Errorf("BD%v @ %v,%v: "+f, args...)
			if !failed {
				failed = true
				failCount++
			}
		}
		verify := func(what string, got, want any) {
			if !reflect.DeepEqual(got, want) {
				errorf("%s:\nwant: %v\n got: %v", what, want, got)
			}
		}

		src := td.Packed(x, y, bd.Planes())
		got := td.BaseImage()

		c := tileconv.Packed{BitDepth: bd}
		c.Decode(src, got, x, y)

		goodSrc := td.Packed(x, y, bd.Planes())
		verify("source data was corrupted", src, goodSrc)

		verify("Stride was changed", got.Stride, baseImage.Stride)
		verify("Rect was changed", got.Rect, baseImage.Rect)
		verify("Palette was changed", got.Palette, baseImage.Palette)

		badPixel := func(c int, m string, tx, ty int, w, g uint8) {
			if c < maxBadPix {
				errorf("%s at %v,%v: want %v, got %v", m, tx, ty, w, g)
			}
		}

		inCount, outCount := 0, 0
		b := baseImage.Bounds()
		area := image.Rect(x, y, x+8, y+8)
		mask := bd.ColorMask()
		for ty := b.Min.Y; ty < b.Max.Y; ty++ {
			for tx := b.Min.X; tx < b.Max.X; tx++ {
				g := got.ColorIndexAt(tx, ty)
				var w uint8

				inside := image.Point{tx, ty}.In(area)
				if inside {
					w = mask & fullImage.ColorIndexAt(tx, ty)
				} else {
					w = baseImage.ColorIndexAt(tx, ty)
				}
				if g == w {
					continue
				}
				if inside {
					badPixel(inCount, "bad pixel", tx, ty, w, g)
					inCount++
				} else {
					badPixel(outCount, "pixel corrupted", tx, ty, w, g)
					outCount++
				}
			}
		}
		if inCount > 0 {
			errorf("%v bad pixels in decode area", inCount)
		}
		if outCount > 0 {
			errorf("%v corrupted pixels outside decode area", outCount)
		}
	}

	runTestForDepth := func(t *testing.T, bd tileconv.BitDepth) {
		failCount = 0
		for y := -4; y < 4 && failCount < maxFails; y++ {
			for x := -4; x < 4 && failCount < maxFails; x++ {
				check(t, bd, x, y)
			}
		}
		if failCount >= maxFails {
			t.Logf("too many errors, skipping rest of depth %v", bd)
		}
	}

	names := []string{"Seq", "Rng"}
	for i, rng := range []bool{false, true} {
		t.Run(names[i], func(t *testing.T) {
			td = newTestImageData(rng)
			fullImage = td.FullImage()
			baseImage = td.BaseImage()

			for bd := tileconv.BD1; bd <= tileconv.BD8; bd++ {
				runTestForDepth(t, bd)
			}
		})
	}
}
