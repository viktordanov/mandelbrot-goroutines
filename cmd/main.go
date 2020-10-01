package main

import (
	"image"
	"image/png"
	"mandelbrot/mandelbrot"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	dimensions := mandelbrot.Size{X: 12000, Y: 8000}
	canvas := make([]uint8, int(dimensions.X*dimensions.Y)*4)
	mandelbrot.Mandelbrot(canvas, dimensions, 2000)
	err := writeBytesToImage(canvas, dimensions)
	if err != nil {
		panic(err)
	}
}

func writeBytesToImage(data []uint8, size mandelbrot.Size) error {
	// Create a blank image 100x200 pixels
	myImage := image.NewRGBA(image.Rect(0, 0, int(size.X), int(size.Y)))
	myImage.Pix = data
	// outputFile is a File type which satisfies Writer interface
	outputFile, err := os.Create("test.png")
	if err != nil {
		// Handle error
		return err
	}

	// Don't forget to close files
	defer outputFile.Close()

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	return png.Encode(outputFile, myImage)
}
