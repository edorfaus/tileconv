package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/png"
	"os"

	"github.com/alexflint/go-arg"

	"github.com/edorfaus/tileconv"
)

func main() {
	var args Args
	arg.MustParse(&args)
	if err := run(args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

type Args struct {
	Input  string `arg:"positional,required" help:"input image file"`
	Output string `arg:"positional,required" help:"output file"`
	Format Format `arg:"-f,required" help:"tile data format; see below"`

	Bpp tileconv.BitDepth `arg:"-b,required" help:"bits per pixel; 1-8"`
}

type Format string

func (Args) Epilogue() string {
	return `Tile data formats:
    p, packed               : packed-pixel
    tp, tileplanar          : planar, per tile
    rp, rowplanar           : planar, per row
    trpp, tilerowpairplanar : planar, pairs per row, rest per tile`
}

func (f *Format) UnmarshalText(text []byte) error {
	switch string(text) {
	case "p", "packed":
	case "tp", "tileplanar":
	case "rp", "rowplanar":
	case "trpp", "tilerowpairplanar":
	default:
		return fmt.Errorf("unknown tile format %q", text)
	}
	*f = Format(text)
	return nil
}

func run(args Args) (e error) {
	var codec tileconv.Codec

	switch args.Format {
	case "p", "packed":
		codec = tileconv.Packed{
			BitDepth: args.Bpp,
		}
	case "tp", "tileplanar":
		codec = tileconv.TilePlanar{
			BitDepth: args.Bpp,
		}
	case "rp", "rowplanar":
		codec = tileconv.RowPlanar{
			BitDepth: args.Bpp,
		}
	case "trpp", "tilerowpairplanar":
		codec = tileconv.TileRowPairPlanar{
			BitDepth: args.Bpp,
		}
	default:
		return fmt.Errorf("unknown tile format: %q", args.Format)
	}

	img, err := loadImage(args.Input)
	if err != nil {
		return err
	}

	out, err := os.Create(args.Output)
	if err != nil {
		return err
	}
	defer tailError(&e, out.Close)

	if err := tileconv.Encode(img, out, codec); err != nil {
		return err
	}

	return nil
}

func loadImage(fn string) (_ image.PalettedImage, e error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer tailError(&e, f.Close)

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	if pi, ok := img.(image.PalettedImage); ok {
		return pi, nil
	}

	return nil, fmt.Errorf("not a paletted image: %s", fn)
}

func tailError(err *error, fn func() error) {
	if e := fn(); e != nil && err != nil && *err == nil {
		*err = e
	}
}
