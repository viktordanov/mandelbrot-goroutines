package main

import (
	"image"
	"image/png"
	"mandelbrot/pkg"
	"os"
	"runtime"
)

func char(n byte) string {
	if float64(n) < pkg.Threshold/2 {
		return " "
	}
	return "*"
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	dimensions := pkg.Size{X: 12000, Y: 8000}
	canvas := make([]uint8, int(dimensions.X*dimensions.Y)*4)
	pkg.Mandelbrot(canvas, dimensions, 2000)
	writeBytesToImage(canvas, dimensions)
}

func writeBytesToImage(data []uint8, size pkg.Size) {
	// Create a blank image 100x200 pixels
	myImage := image.NewRGBA(image.Rect(0, 0, int(size.X), int(size.Y)))
	myImage.Pix = data
	// outputFile is a File type which satisfies Writer interface
	outputFile, err := os.Create("test.png")
	if err != nil {
		// Handle error
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(outputFile, myImage)

	// Don't forget to close files
	outputFile.Close()
}
