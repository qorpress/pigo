package pigo

import (
	"fmt"
	"image"
	"image/color"

	"github.com/esimov/stackblur-go"
)

// SubImager is a wrapper implementing the SubImage method from the image package.
type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

// BlinkDetector detects blink occurrence.
func BlinkDetector(puploc *Puploc, img image.Image) (image.Image, bool, error) {
	rect := image.Rect(
		puploc.Col-int(puploc.Scale*1.5),
		puploc.Row-int(puploc.Scale*1.5),
		puploc.Col+int(puploc.Scale*1.5),
		puploc.Row+int(puploc.Scale*1.5),
	)
	//fmt.Println(rect.Bounds())
	subImg := img.(SubImager).SubImage(rect)
	dx, dy := subImg.Bounds().Dx(), subImg.Bounds().Dy()

	blur, err := stackblur.Run(rgbToGrayscale(subImg.(*image.NRGBA)), 5)
	if err != nil {
		return nil, false, err
	}

	res := Sobel(blur.(*image.NRGBA), 70)

	// Use pixel thresholding to obtain a full black and white sub image.
	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			pixel := res.At(x, y)
			threshold := func(pixel color.Color) color.Color {
				r, _, _, _ := pixel.RGBA()
				if r > 127 {
					return color.RGBA{R: 255, G: 255, B: 255, A: 255}
				}
				return color.RGBA{R: 0, G: 0, B: 0, A: 255}
			}
			res.Set(x, y, threshold(pixel))
		}
	}
	if ok := detectBlink(res, 0.44); !ok {
		fmt.Println("Blink detected")
		return res, true, nil
	}
	return res, false, nil
}

func detectBlink(src *image.NRGBA, th float64) bool {
	bounds := src.Bounds()
	//fmt.Println(int(float64(bounds.Max.X) * 0.05))
	//fmt.Println(bounds)
	startX, startY := bounds.Max.X/2, bounds.Max.Y/2
	cx1, cx2, cy1, cy2, ratio := 0, 0, 0, 0, 0.0

	// Traverse each pixel to the far right of the eye bounding box
	// starting from the pupil's center point. Break if a white pixel is reached.
	// Since the Sobel operator will give us the contours of the eye/pupil, we can localize the pupil margins this way.
	for x := startX; x <= bounds.Max.X; x++ {
		r, g, b, _ := src.At(x, startY).RGBA()
		if r>>8 == 0 {
			src.Set(x, startY, color.RGBA{R: 255, G: uint8(g >> 8), B: uint8(b >> 8), A: 255})
		}
		if r>>8 == 255 {
			break
		}
		cx1++
	}

	// Traverse each pixel to the far left of the eye bounding box
	// starting from the pupil's center point. Break if a white pixel is reached.
	for x := startX - 1; x >= bounds.Min.X; x-- {
		r, g, b, _ := src.At(x, startY).RGBA()
		if r>>8 == 0 {
			src.Set(x, startY, color.RGBA{R: 255, G: uint8(g >> 8), B: uint8(b >> 8), A: 255})
		}
		if r>>8 == 255 {
			break
		}
		cx2++
	}

	// Traverse each pixel to the bottom of the eye bounding box
	// starting from the pupil's center point. Break if a white pixel is reached.
	for y := startY + 1; y <= bounds.Max.Y; y++ {
		r, g, b, _ := src.At(startX, y).RGBA()
		if r>>8 == 0 {
			src.Set(startX, y, color.RGBA{R: 255, G: uint8(g >> 8), B: uint8(b >> 8), A: 255})
		}
		if r>>8 == 255 {
			break
		}
		cy1++
	}

	// Traverse each pixel to the top of the eye bounding box
	// starting from the pupil's center point. Break if a white pixel is reached.
	for y := startY - 1; y >= bounds.Min.Y; y-- {
		r, g, b, _ := src.At(startX, y).RGBA()
		if r>>8 == 0 {
			src.Set(startX, y, color.RGBA{R: 255, G: uint8(g >> 8), B: uint8(b >> 8), A: 255})
		}
		if r>>8 == 255 {
			break
		}
		cy2++
	}

	// Check if we reached the eyes bounding boxes. This means that the eyes are closed.
	// The reason why we are checking the bounding box - 1 is that the
	// sobel operator is marking the last row and column as white.
	if cx1 == bounds.Max.X/2-1 {
		cx1 = 0
	} else if cx2 == bounds.Max.X/2 {
		cx2 = 0
	} else if cy1 == bounds.Max.Y/2-1 {
		cy1 = 0
	} else if cy2 == bounds.Max.Y/2 {
		cy2 = 0
	}
	// if cx1 == 0 {
	// 	cx2 = 0
	// } else if cx2 == 0 {
	// 	cx1 = 0
	// }

	cx := cx1 + cx2
	cy := cy1 + cy2
	if cx == 0 { // we do not need to check vertically.
		cy = 0
	}
	if cx > 0 && cy > 0 {
		ratio = float64(cy) / float64(cx)
	}
	fmt.Println("XY:", cx, cy)
	fmt.Println("Ratio:", ratio)
	if ratio > th {
		return true
	}
	return false
}
