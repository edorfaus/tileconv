package chrconv_test

import (
	"testing"

	"github.com/edorfaus/chrconv"
)

func TestBytesPerPlane(t *testing.T) {
	// NOTE: if this value is changed for any reason, that breaks a lot
	// of assumptions elsewhere in the code - so be sure to update that
	// code as well, not just this test.
	if chrconv.BytesPerPlane != 8 {
		t.Errorf("wrong value: want 8, got %v", chrconv.BytesPerPlane)
	}
}

func TestBitDepthPlanes(t *testing.T) {
	check := func(bd chrconv.BitDepth, want int) {
		t.Helper()
		got := bd.Planes()
		if got != want {
			t.Errorf("depth %v planes: want %v, got %v", bd, want, got)
		}
	}
	check(chrconv.BD1, 1)
	check(chrconv.BD2, 2)
	check(chrconv.BD3, 3)
	check(chrconv.BD4, 4)
	check(chrconv.BD5, 5)
	check(chrconv.BD6, 6)
	check(chrconv.BD7, 7)
	check(chrconv.BD8, 8)
}

func TestBitDepthColors(t *testing.T) {
	check := func(bd chrconv.BitDepth, want int) {
		t.Helper()
		got := bd.Colors()
		if got != want {
			t.Errorf("depth %v colors: want %v, got %v", bd, want, got)
		}
	}
	check(chrconv.BD1, 2)
	check(chrconv.BD2, 4)
	check(chrconv.BD3, 8)
	check(chrconv.BD4, 16)
	check(chrconv.BD5, 32)
	check(chrconv.BD6, 64)
	check(chrconv.BD7, 128)
	check(chrconv.BD8, 256)
}

func TestBitDepthColorMask(t *testing.T) {
	check := func(d chrconv.BitDepth, want uint8) {
		t.Helper()
		got := d.ColorMask()
		if got != want {
			t.Errorf("depth %v mask: want %08b, got %08b", d, want, got)
		}
	}
	check(chrconv.BD1, 0b00000001)
	check(chrconv.BD2, 0b00000011)
	check(chrconv.BD3, 0b00000111)
	check(chrconv.BD4, 0b00001111)
	check(chrconv.BD5, 0b00011111)
	check(chrconv.BD6, 0b00111111)
	check(chrconv.BD7, 0b01111111)
	check(chrconv.BD8, 0b11111111)
}

func TestBitDepthBytesPerTile(t *testing.T) {
	check := func(bd chrconv.BitDepth, want int) {
		t.Helper()
		got := bd.BytesPerTile()
		if got != want {
			t.Errorf("depth %v bytes: want %v, got %v", bd, want, got)
		}
	}
	check(chrconv.BD1, 8*1)
	check(chrconv.BD2, 8*2)
	check(chrconv.BD3, 8*3)
	check(chrconv.BD4, 8*4)
	check(chrconv.BD5, 8*5)
	check(chrconv.BD6, 8*6)
	check(chrconv.BD7, 8*7)
	check(chrconv.BD8, 8*8)
}
