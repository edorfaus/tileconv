package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
	"strings"

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
	Input  string `arg:"positional,required" help:"input file"`
	Output string `arg:"positional,required" help:"output file"`
	Decode bool   `arg:"-d" help:"decode tiles into an image"`
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

	if args.Decode {
		return runDecode(args, codec)
	}

	return runEncode(args, codec)
}

func runEncode(args Args, codec tileconv.Codec) (e error) {
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

func runDecode(args Args, codec tileconv.Codec) (e error) {
	outFmt := strings.ToLower(filepath.Ext(args.Output))
	if outFmt != ".png" && outFmt != ".gif" {
		return fmt.Errorf("unknown image format: %q", outFmt)
	}

	src, err := os.ReadFile(args.Input)
	if err != nil {
		return err
	}

	if len(src)%codec.Size() != 0 {
		return fmt.Errorf("input is not a whole number of tiles")
	}

	tiles := len(src) / codec.Size()
	rows := (tiles + 15) / 16
	cols := 16
	if rows < 2 {
		cols = tiles
	}

	img := image.NewPaletted(
		image.Rect(0, 0, cols*8, rows*8), makePalette(args.Bpp),
	)

	tileconv.Decode(src, img, codec)

	out, err := os.Create(args.Output)
	if err != nil {
		return err
	}
	defer tailError(&e, out.Close)

	switch outFmt {
	case ".png":
		return png.Encode(out, img)
	case ".gif":
		// When given an *image.Paletted (like we're doing), and the
		// options.NumColors matches its palette, then the Quantizer and
		// Drawer are not actually used, and the image is used as-is.
		return gif.Encode(out, img, &gif.Options{
			NumColors: len(img.Palette),
		})
	default:
		return fmt.Errorf("unexpected image format: %q", outFmt)
	}
}

func makePalette(bpp tileconv.BitDepth) color.Palette {
	p := make(color.Palette, bpp.Colors())
	max := len(p) - 1
	for i := 0; i < len(p); i++ {
		v := uint8(255 * i / max)
		p[i] = color.NRGBA{
			R: v, G: v, B: v, A: 255,
		}
	}
	return p
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
