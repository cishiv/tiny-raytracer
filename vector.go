package main

import "math"

// Vec3 : 3 dimensional float64 vector representation
type Vec3 struct {
	x float64
	y float64
	z float64
}

// NewVec3 : Construct and return a new Vec3 instance
func NewVec3(x, y, z float64) Vec3 {
	return Vec3{
		x,
		y,
		z,
	}
}

// MultiplyScalar : Vec3 Scalar Multiplication
func (v Vec3) MultiplyScalar(scalar float64) Vec3 {
	return Vec3{
		v.x * scalar,
		v.y * scalar,
		v.z * scalar,
	}
}

// ConvertToRBG : Return a triplet of unsigned 8-bit ints corresponding to rgb values from a vector v
func (v Vec3) ConvertToRGB() (uint8, uint8, uint8) {
	return uint8(255 * math.Max(0.0, math.Min(1.0, v.x))),
		uint8(255 * math.Max(0.0, math.Min(1.0, v.y))),
		uint8(255 * math.Max(0.0, math.Min(1.0, v.z)))
}
