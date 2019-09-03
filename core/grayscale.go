package pigo

import (
	"image"
	"image/color"
)

// RgbToGrayscale converts the image to grayscale mode and returns as an uint array.
func RgbToGrayscale(src image.Image) []uint8 {
	cols, rows := src.Bounds().Dx(), src.Bounds().Dy()
	gray := make([]uint8, rows*cols)

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			gray[y*cols+x] = uint8(
				(0.299*float64(r) +
					0.587*float64(g) +
					0.114*float64(b)) / 256,
			)
		}
	}
	return gray
}

// rgbToGrayscale converts the image to grayscale mode but it returns as an image.
func rgbToGrayscale(src *image.NRGBA) *image.NRGBA {
	dx, dy := src.Bounds().Max.X, src.Bounds().Max.Y
	dst := image.NewNRGBA(src.Bounds())
	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			r, g, b, _ := src.At(x, y).RGBA()
			lum := float32(r)*0.299 + float32(g)*0.587 + float32(b)*0.114
			pixel := color.Gray{Y: uint8(lum / 256)}
			dst.Set(x, y, pixel)
		}
	}
	return dst
}
