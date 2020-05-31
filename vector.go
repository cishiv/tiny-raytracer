package main

import "math"

// --- VECTOR STRUCTS AND UTILITIES --
// Author: https://github.com/cishiv

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

// Vec2 : 2 dimensional f64 vector representation
type Vec2 struct {
	x float64
	y float64
}

// NewVec2 : Construct and return a new Vec2 instance
func NewVec2(x, y float64) Vec2 {
	return Vec2{
		x,
		y,
	}
}

// Vec4 : Vec3 + w, '4 dimensional vector'
type Vec4 struct {
	x float64
	y float64
	z float64
	w float64
}

// NewVec4 : Construct and return a new Vec4 instance
func NewVec4(x, y, z, w float64) Vec4 {
	return Vec4{
		x,
		y,
		z,
		w,
	}
}

// DotProduct : Give Vec3 (u) and (v) return the dot production u.DotProduct(v)
func (u Vec3) DotProduct(v Vec3) float64 {
	return u.x*v.x + u.y*v.y + u.z*v.z
}

// Magnitude : Return the magnitude of a Vec3 u
func (u Vec3) Magnitude() float64 {
	return math.Sqrt(u.DotProduct(u))
}

// CrossProduct : Return u x v (Vec3)
func (u Vec3) CrossProduct(v Vec3) Vec3 {
	return Vec3{
		u.y*v.z - u.z*v.y,
		u.z*v.x - u.x*v.z,
		u.x*v.y - u.y*v.x,
	}
}

// Normalize : Return a normalized Vec3 given u
func (u Vec3) Normalize() Vec3 {
	length := u.Magnitude()
	if length > 0 {
		return u.MultiplyScalar(1.0 / length)
	}
	return u
}

// Add : Given u and v, return u + v
func (u Vec3) Add(v Vec3) Vec3 {
	return Vec3{
		u.x + v.x,
		u.y + v.y,
		u.z + v.z,
	}
}

// Subtract : Given u and v, return u - v
func (u Vec3) Subtract(v Vec3) Vec3 {
	return Vec3{
		u.x - v.x,
		u.y - v.y,
		u.z - v.z,
	}
}

// Negate : Return a negated Vec3 given u
func (u Vec3) Negate() Vec3 {
	return Vec3{
		-u.x,
		-u.y,
		-u.z,
	}
}

// MultiplyScalar : Vec3 Scalar Multiplication
func (u Vec3) MultiplyScalar(scalar float64) Vec3 {
	return Vec3{
		u.x * scalar,
		u.y * scalar,
		u.z * scalar,
	}
}

// Copy : Copy values from v into u
func (u *Vec3) Copy(v Vec3) {
	u.x = v.x
	u.y = v.y
	u.z = v.z
}

// ConvertToRGB : Return a triplet of unsigned 8-bit ints corresponding to rgb values from a vector u
func (u Vec3) ConvertToRGB() (uint8, uint8, uint8) {
	return uint8(255 * math.Max(0.0, math.Min(1.0, u.x))),
		uint8(255 * math.Max(0.0, math.Min(1.0, u.y))),
		uint8(255 * math.Max(0.0, math.Min(1.0, u.z)))
}
