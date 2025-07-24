/*
Package tileconv provides conversion between images and several retro
8x8 tile graphics formats.

It is targeted at consoles and images that use indexed-color (paletted)
graphics, and at images that are specifically designed for that console.

As such, this package works on the color indexes themselves, rather than
their colors, since that is more useful and easier to reason about.

# Organization

The main abstraction that this package provides is that of the [Codec]
interface. A Codec represents a tile graphics format, and provides a way
to encode or decode a single tile using that format.

Additional functionality, like handling multiple tiles, is then built on
top of that abstraction.

This makes it fairly easy both to pick which format you need to use, and
to extend the library with other formats if necessary.

To avoid a proliferation of types, it is common for an implementation to
support multiple closely related formats, with a field to specify which.

This is where the [BitDepth] type comes in, to specify which bit depth
(how many bits per pixel) should be used during the conversion.

Note that the codec implementations in this package often support more
variations (e.g. bit depths) than are supported by the retro consoles
themselves, so you still need to do your own due diligence on that.
*/
package tileconv
