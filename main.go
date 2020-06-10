package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// Author: https://github.com/cishiv

// Inspired by https://github.com/ssloy/tinyraytracer
// https://github.com/huqa/tinyraytracer-go/blob/54817bbe85026965ad4d0725a1375705f82d8e09/internal/app/tinyraytracer-go/render.go#L110
// https://github.com/ssloy/tinyraytracer/blob/bd36c9857305b3cbd06f5b768bb48a92df9ae68b/geometry.h
// https://github.com/huqa/tinyraytracer-go/blob/master/internal/pkg/vector/vector.go

// Camera view angle
const fieldOfView = math.Pi / 3.0

// 'Window' dimensions (image dimensions since we're writing directly to file)
const width = 1024
const height = 768

// Sample materials
var ivory = NewMaterial(NewVec3(0.4, 0.4, 0.3), NewVec4(0.6, 0.3, 0.1, 0.0), 50, 1.0)
var redRubber = NewMaterial(NewVec3(0.3, 0.1, 0.1), NewVec4(0.9, 0.1, 0.1, 0.0), 10, 1.0)
var mirror = NewMaterial(NewVec3(1.0, 1.0, 1.0), NewVec4(0.0, 10.0, 0.8, 0.0), 1425, 1.0)
var glass = NewMaterial(NewVec3(0.6, 0.7, 0.8), NewVec4(0.0, 0.5, 0.1, 0.8), 125, 1.5)

func render(spheres []Sphere, lights []Light) {
	// Construct the framebuffer
	framebuffer := make([]Vec3, width*height)
	origin := NewVec3(0, 0, 0)
	// Fill the framebuffer
	fmt.Println("Filling framebuffer")
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			// Big math (TODO: Figure this out)
			x := (float64(i) + 0.5) - width/2.0
			y := -(float64(j) + 0.5) + height/2.0
			// TODO: We're not 3D yet
			z := -float64(height) / (2.0 * math.Tan(fieldOfView/2.0))
			direction := NewVec3(x, y, z).Normalize()
			framebuffer[i+j*width] = CastRay(&origin, &direction, spheres, lights, 0)
		}
	}
	fmt.Println("Framebuffer filled successfully")

	// Save to file
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	fmt.Println("saving framebuffer to file")

	// Iterate over the framebuffer and fill the image based on the gradient corresponding to the vectors in the buffer
	var max float64
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			v := framebuffer[i+j*width]
			max = math.Max(v.x, math.Max(v.y, v.z))
			if max > 1 {
				v = v.MultiplyScalar((1.0 / max))
			}
			r, g, b := v.ConvertToRGB()
			c := color.RGBA{r, g, b, 255}
			img.Set(i, j, c)
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
	render(createSpheres(), createLights())
}

func createSpheres() []Sphere {
	spheres := make([]Sphere, 0)
	spheres = append(spheres, NewSphere(-3, 0, -16, 2, ivory))
	spheres = append(spheres, NewSphere(-1.0, -1.5, -12, 2, glass))
	spheres = append(spheres, NewSphere(1.5, -0.5, -18, 3, redRubber))
	spheres = append(spheres, NewSphere(7, 5, -18, 4, mirror))
	fmt.Println(spheres)
	return spheres
}

func createLights() []Light {
	lights := make([]Light, 0)
	lights = append(lights, NewLight(-20, 20, 20, 4.5))
	lights = append(lights, NewLight(20, -20, -20, 1.8))
	lights = append(lights, NewLight(30, 20, 30, 1.7))
	return lights
}
