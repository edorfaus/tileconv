package tileconv_test

import (
	"testing"

	"github.com/edorfaus/tileconv"
)

func TestBytesPerPlane(t *testing.T) {
	// NOTE: if this value is changed for any reason, that breaks a lot
	// of assumptions elsewhere in the code - so be sure to update that
	// code as well, not just this test.
	if tileconv.BytesPerPlane != 8 {
		t.Errorf("wrong value: want 8, got %v", tileconv.BytesPerPlane)
	}
}

func TestBitDepthPlanes(t *testing.T) {
	check := func(bd tileconv.BitDepth, want int) {
		t.Helper()
		got := bd.Planes()
		if got != want {
			t.Errorf("depth %v planes: want %v, got %v", bd, want, got)
		}
	}
	check(tileconv.BD1, 1)
	check(tileconv.BD2, 2)
	check(tileconv.BD3, 3)
	check(tileconv.BD4, 4)
	check(tileconv.BD5, 5)
	check(tileconv.BD6, 6)
	check(tileconv.BD7, 7)
	check(tileconv.BD8, 8)
}

func TestBitDepthColors(t *testing.T) {
	check := func(bd tileconv.BitDepth, want int) {
		t.Helper()
		got := bd.Colors()
		if got != want {
			t.Errorf("depth %v colors: want %v, got %v", bd, want, got)
		}
	}
	check(tileconv.BD1, 2)
	check(tileconv.BD2, 4)
	check(tileconv.BD3, 8)
	check(tileconv.BD4, 16)
	check(tileconv.BD5, 32)
	check(tileconv.BD6, 64)
	check(tileconv.BD7, 128)
	check(tileconv.BD8, 256)
}

func TestBitDepthColorMask(t *testing.T) {
	check := func(d tileconv.BitDepth, want uint8) {
		t.Helper()
		got := d.ColorMask()
		if got != want {
			t.Errorf("depth %v mask: want %08b, got %08b", d, want, got)
		}
	}
	check(tileconv.BD1, 0b00000001)
	check(tileconv.BD2, 0b00000011)
	check(tileconv.BD3, 0b00000111)
	check(tileconv.BD4, 0b00001111)
	check(tileconv.BD5, 0b00011111)
	check(tileconv.BD6, 0b00111111)
	check(tileconv.BD7, 0b01111111)
	check(tileconv.BD8, 0b11111111)
}

func TestBitDepthBytesPerTile(t *testing.T) {
	check := func(bd tileconv.BitDepth, want int) {
		t.Helper()
		got := bd.BytesPerTile()
		if got != want {
			t.Errorf("depth %v bytes: want %v, got %v", bd, want, got)
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
