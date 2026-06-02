package icons

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

// indigo-500 from the Tailwind palette — matches the frontend theme
var fill = color.RGBA{R: 99, G: 102, B: 241, A: 255}

// Icon returns the bytes of a 32×32 PNG used as the systray icon.
// Replace with a proper icon file (via go:embed) when assets are ready.
func Icon() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))

	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			// 4-px padding on all sides gives a 24×24 filled square
			if x >= 4 && x < 28 && y >= 4 && y < 28 {
				img.Set(x, y, fill)
			}
		}
	}

	var buf bytes.Buffer
	png.Encode(&buf, img) //nolint:errcheck — in-memory encode never fails
	return buf.Bytes()
}
