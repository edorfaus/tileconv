TileConv
========

[![Go Reference][GoRefImage]][GoRefLink]

[GoRefImage]: https://pkg.go.dev/badge/github.com/edorfaus/tileconv.svg
[GoRefLink]: https://pkg.go.dev/github.com/edorfaus/tileconv

This project is aimed at being a library for converting images to retro
8x8 tile graphics formats (and back), and for convenience also includes
a CLI tool that uses the library to do that conversion for basic images.

It is targeted at consoles and images that use indexed-color (paletted)
graphics, and at images that are specifically designed for that console.

As such, this library works on the color indexes themselves, rather than
their colors, since that is more useful and easier to reason about.

Each kind of tile graphics format is handled by an implementation of the
Codec interface, usually with a field to specify the bit depth you want.

E.g. for a basic packed-pixel format, with 4 bits per pixel:

```go
	codec := tileconv.Packed{BitDepth: tileconv.BD4}
```

Some additional functionality, like handling multiple tiles per image,
is built on top of that interface as separate top-level functions.

E.g. to convert an image with multiple tiles to an output file:

```go
	err := tileconv.Encode(img, outfile, codec)
```

Note that the codec implementations in this package often support more
variations (e.g. bit depths) than are supported by the retro consoles
themselves, so you still need to do your own due diligence on that.
