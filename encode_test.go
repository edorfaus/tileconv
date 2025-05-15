package chrconv_test

import (
	"bytes"
	"image"
	"reflect"
	"testing"

	"github.com/edorfaus/chrconv"
)

// TestEncode_Data tests that Encode() actually calls the codec properly
// and writes the returned data to the Writer in the expected order.
func TestEncode_Data(t *testing.T) {
	var td *testImageData
	var src, goodSrc *image.Paletted

	runTest := func(t *testing.T, copies int) {
		c := &testCodec{t: t, copies: copies}
		w := &bytes.Buffer{}

		err := chrconv.Encode(src, w, c)

		if err != nil {
			t.Errorf("unexpected encode error (%T): %v", err, err)
		}

		if !reflect.DeepEqual(src, goodSrc) {
			t.Errorf("source image was corrupted")
		}

		verify := func(msg string, want, got int) {
			t.Helper()
			if got != want {
				t.Errorf("bad %s: want %v, got %v", msg, want, got)
			}
		}

		verify("encode call count", 4, c.encodes)
		verify("decode call count", 0, c.decodes)

		got := w.Bytes()

		sz := c.Size()
		verify("data length", 4*sz, len(got))

		want := make([]byte, sz)
		checkAt := func(x, y, chunk int) {
			t.Helper()
			var g []byte
			ofs := chunk * sz
			if ofs >= len(got) {
				t.Errorf("chunk %v (%v,%v): got no data", chunk, x, y)
				return
			}
			if len(got) >= ofs+sz {
				g = got[ofs : ofs+sz]
			} else {
				g = got[ofs:]
				t.Errorf(
					"chunk %v (%v,%v): got only %v of %v bytes",
					chunk, x, y, len(g), sz,
				)
			}

			c.Encode(src, x, y, want)
			// Since different lengths already caused a fail above, this
			// only checks for the values being different, up to length.
			if !reflect.DeepEqual(g, want[:len(g)]) {
				t.Errorf(
					"chunk %v (%v,%v): bad data:\nwant: %v\n got: %v",
					chunk, x, y, want, g,
				)
			}
		}

		checkAt(-4, -4, 0)
		checkAt(4, -4, 1)
		checkAt(-4, 4, 2)
		checkAt(4, 4, 3)
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

			t.Run("small", func(t *testing.T) {
				src = td.FullImage()
				runTest(t, 1)
			})
			t.Run("large", func(t *testing.T) {
				src = td.FullImage()
				runTest(t, 2)
			})
			// An image with a size that is not an even multiple of the
			// tile size shall be encoded as if the size was rounded up.
			t.Run("non-even", func(t *testing.T) {
				area := image.Rect(-4, -4, 5, 5)
				goodSrc = goodSrc.SubImage(area).(*image.Paletted)
				src = td.FullImage().SubImage(area).(*image.Paletted)
				runTest(t, 1)
			})
		})
	}
}

// TestEncode_Error tests that Encode() returns the errors that it got
// from the Writer it writes to.
func TestEncode_Error(t *testing.T) {
	td := newTestImageData(false)
	src := td.FullImage()
	c := &testCodec{t: t, copies: 1}

	testFn := func(remain int, expected bool) func(*testing.T) {
		return func(t *testing.T) {
			t.Helper()
			w := &ErrWriter{Remain: remain}

			err := chrconv.Encode(src, w, c)

			if err == nil {
				if expected {
					t.Errorf("missing error (%v remain)", w.Remain)
				}
			} else if err != w {
				t.Errorf(
					"wrong error:\nwant: %#v\n got: %#v : %v",
					w, err, err,
				)
			} else if !expected {
				t.Errorf(
					"unexpected error (%v remain): %#v : %v",
					w.Remain, err, err,
				)
			}
		}
	}

	t.Run("immediate", testFn(0, true))
	t.Run("halfway", testFn(2*c.Size(), true))
	t.Run("at end", testFn(4*c.Size(), true))
	t.Run("none", testFn(4*c.Size()+1, false))
}

// ErrWriter is an io.Writer that returns an error (itself) as soon as a
// certain number of bytes have been written to it.
type ErrWriter struct {
	Remain int
}

// Write implements io.Writer.
func (w *ErrWriter) Write(b []byte) (int, error) {
	if len(b) < w.Remain {
		w.Remain -= len(b)
		return len(b), nil
	}
	rem := w.Remain
	w.Remain = 0
	return rem, w
}

// Error implements error.
func (w *ErrWriter) Error() string {
	return "ErrWriter"
}
