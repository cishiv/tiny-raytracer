package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// Inspired by https://github.com/ssloy/tinyraytracer

func render() {
	width := 1024
	height := 768

	// Construct the framebuffer
	framebuffer := make([]Vec3, width*height)

	// Fill the framebuffer
	fmt.Println("Filling framebuffer")
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			framebuffer[j+i*width] = NewVec3(float64(i)/float64(height), float64(j)/float64(width), 0)
		}
	}
	fmt.Println("Framebuffer filled successfully")

	// Save to file
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	fmt.Println("saving framebuffer to file")

	// Iterate over the framebuffer and fill the image based on the gradient corresponding to the vectors in the buffer
	var max float64
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			v := framebuffer[j+i*width]
			max = math.Max(v.x, math.Max(v.y, v.z))
			if max > 1 {
				v = v.MultiplyScalar((1.0 / max))
			}
			r, g, b := v.ConvertToRGB()
			c := color.RGBA{r, g, b, 255}
			img.Set(j, i, c)
		}
	}

	// Write the file
	out, err := os.Create("test.png")
	if err != nil {
		fmt.Println("error occurred creating file")
	}
	defer out.Close()
	png.Encode(out, img)
}

func main() {
	render()
}
