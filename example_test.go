package tileconv_test

import (
	"fmt"
	"image"
	"image/color/palette"
	"os"

	"github.com/edorfaus/tileconv"
)

func Example() {
	input := []byte("12345678abcdefgh")

	tile := image.NewPaletted(image.Rect(0, 0, 8, 8), palette.Plan9)

	codec1 := tileconv.TilePlanar{BitDepth: tileconv.BD2}
	tileconv.Decode(input, tile, codec1)

	codec2 := tileconv.RowPlanar{BitDepth: tileconv.BD2}
	err := tileconv.Encode(tile, os.Stdout, codec2)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}

	// Output: 1a2b3c4d5e6f7g8h
}
