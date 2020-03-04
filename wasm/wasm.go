// +build js,wasm

package main

import (
	"github.com/qorpress/pigo/wasm/canvas"
)

func main() {
	c := canvas.NewCanvas()
	webcam, err := c.StartWebcam()
	if err != nil {
		c.Alert("Webcam not detected!")
	} else {
		webcam.Render()
	}
}
