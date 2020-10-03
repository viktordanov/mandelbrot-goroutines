package mandelbrot

import (
	"fmt"
	"math"
	"math/cmplx"
	"sync"
)

// Size defines the size.
type Size struct {
	X float64
	Y float64
}

// Global config consts.
const (
	Threshold float64 = 2
	ReStart   float64 = -2.0
	ReEnd     float64 = 1.0
	ImStart   float64 = -1.0
	ImEnd     float64 = 1.0
	OffsetX   float64 = 750.25 / 6000.0
	OffsetY   float64 = 1685.1851 / 4000.0
	Scale     float64 = 0.25
)

// Instructions contains all the instructions.
type Instructions struct {
	r          int
	rMax       int
	i          int
	iMax       int
	size       Size
	iterations int
}

// Pixel is a pixel in an image.
type Pixel struct {
	index int
	data  [4]uint8
}

var test float64
var counter float64

func chunkWorker(ins Instructions, out chan<- Pixel, wg *sync.WaitGroup) {
	for r := ins.r; r < ins.rMax; r++ {
		for i := ins.i; i < ins.iMax; i++ {
			var z complex128
			c := complex(ReStart+(OffsetX*(ReEnd-ReStart))+(float64(r)/ins.size.X)*(ReEnd-ReStart)*Scale,
				ImStart+(OffsetY*(ImEnd-ImStart))+(float64(i)/ins.size.Y)*(ImEnd-ImStart)*Scale)
			var n float64
			for n = 0; int(n) < ins.iterations && cmplx.Abs(z) <= Threshold; n++ {
				z = z*z + c
			}

			n -= math.Log(math.Log(cmplx.Abs(z))/math.Log(255)) / math.Log(2.0)
			out <- Pixel{
				r*4 + i*int(ins.size.X*4),
				[4]byte{
					byte(n * 45000 / float64(ins.iterations)),
					0,
					byte(n * 350 / float64(ins.iterations)),
					255,
				},
			}
		}
	}
	fmt.Printf("Progress: %f\n", counter/test)
	counter++
	wg.Done()
}

// Mandelbrot draws the Mandelbrot fractal on the provided canvas.
func Mandelbrot(canvas []uint8, size Size, iterations int) {
	var wg sync.WaitGroup
	out := make(chan Pixel, int(size.X*size.Y))

	chunkSize := 128

	for r := 0; r < int(size.X); r += chunkSize {
		var rMax int
		if r+chunkSize >= int(size.X) {
			rMax = int(size.X)
		} else {
			rMax = r + chunkSize
		}
		for i := 0; i < int(size.Y); i += chunkSize {
			var iMax int
			if i+chunkSize >= int(size.Y) {
				iMax = int(size.Y)
			} else {
				iMax = i + chunkSize
			}
			test++
			wg.Add(1)
			ins := Instructions{r, rMax, i, iMax, size, iterations}
			go chunkWorker(ins, out, &wg)
		}
	}
	wg.Wait()

	for len(out) > 0 {
		pixel := <-out
		canvas[pixel.index] = pixel.data[0]
		canvas[pixel.index+1] = pixel.data[1]
		canvas[pixel.index+2] = pixel.data[2]
		canvas[pixel.index+3] = pixel.data[3]
	}
}
