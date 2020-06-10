package main

import (
	"math"
)

// Author: https://github.com/cishiv

// COLOR
var bgColor = NewVec3(0.2, 0.7, 0.8)
var objectColor = NewVec3(0.4, 0.4, 0.3)

// Sphere : Define a sphere, with a center defined by a 3D vector and a scalar for the radius
type Sphere struct {
	center Vec3
	radius float64
	mat    Material
}

// Material : how do our spheres look?
type Material struct {
	diffuseColor     Vec3
	albedo           Vec3
	specularExponent float64
}

// Light : some basic lighting
type Light struct {
	position  Vec3
	intensity float64
}

// NewLight : Return a new light source from pos(x,y,z) with a particular intensity
func NewLight(x, y, z, intensity float64) Light {
	return Light{
		NewVec3(x, y, z),
		intensity,
	}
}

// NewMaterial : Return a new vec3 representing the diffuse_colour for a material
func NewMaterial(diffuseColor Vec3, albedo Vec3, specularExponent float64) Material {
	return Material{
		diffuseColor,
		albedo,
		specularExponent,
	}
}

// NewSphere : Return a new sphere object
func NewSphere(x, y, z, r float64, mat Material) Sphere {
	return Sphere{
		NewVec3(x, y, z),
		r,
		mat,
	}
}

// SceneIntersect : Cast a ray, determine what spheres are intersected, and what material is cast on
// TODO : Figure SceneIntersect math out
func SceneIntersect(origin, direction, hit, N *Vec3, mat *Material, spheres []Sphere) bool {
	sphereDistance := math.MaxFloat64
	for _, s := range spheres {
		var distI float64
		t0, intersect := s.RayIntersect(*origin, *direction, distI)
		if intersect && t0 < sphereDistance {
			sphereDistance = t0
			k := origin.Add(direction.MultiplyScalar(t0))
			hit.Copy(k)
			n := k.Subtract(s.center).Normalize()
			N.Copy(n)
			mat.diffuseColor = s.mat.diffuseColor
			mat.albedo = s.mat.albedo
			mat.specularExponent = s.mat.specularExponent
		}
	}
	return sphereDistance < 1000.0
}

// CastRay : Cast a ray toward all spheres in a scene, with a given origin of light and a direction, return the resulting 'reflection'(?)
func CastRay(origin, direction *Vec3, spheres []Sphere, lights []Light, depth uint) Vec3 {
	point := &Vec3{}
	N := &Vec3{}
	mat := &Material{}

	if depth > 4 || !SceneIntersect(origin, direction, point, N, mat, spheres) {
		return bgColor
	}

	// LIGHTING
	// ----------

	// offset
	// if the point drawn lines on the surface of the object, any ray will intersect the object (?)
	p := 1E-3

	// Reflection
	// We don't need to normalize here, but just for safety...
	reflectDirection := Reflect(*direction, *N)
	var reflectOrigin Vec3
	if reflectDirection.DotProduct(*N) < 0 {
		reflectOrigin = point.Subtract(N.MultiplyScalar(p))
	} else {
		reflectOrigin = point.Add(N.MultiplyScalar(p))
	}

	reflectColor := CastRay(&reflectOrigin, &reflectDirection, spheres, lights, depth+1)

	diffuseLightIntensity := 0.0
	specularLightIntensity := 0.0
	lightDist := 0.0
	var shadowOrigin Vec3
	for _, light := range lights {
		lightDirection := (light.position.Subtract(*point)).Normalize()
		lightDist = light.position.Subtract(*point).Magnitude()

		if lightDirection.DotProduct(*N) < 0 {
			shadowOrigin = point.Subtract(N.MultiplyScalar(p))
		} else {
			shadowOrigin = point.Add(N.MultiplyScalar(p))
		}

		shadowPoint := Vec3{}
		shadowNormal := Vec3{}
		tempMat := Material{}

		// Shadows: Skip the current light source if there is an intersection on the segment between the light source and current point
		if SceneIntersect(&shadowOrigin, &lightDirection, &shadowPoint, &shadowNormal, &tempMat, spheres) && shadowPoint.Subtract(shadowOrigin).Magnitude() < lightDist {
			continue
		}

		diffuseLightIntensity += light.intensity * math.Max(0.0, lightDirection.DotProduct(*N))

		specularLightIntensity += math.Pow(
			math.Max(0.0, Reflect(lightDirection, *N).DotProduct(*direction)),
			mat.specularExponent,
		) * light.intensity
	}
	lit := mat.diffuseColor.MultiplyScalar(diffuseLightIntensity).MultiplyScalar(mat.albedo.x)
	sp := specularLightIntensity * mat.albedo.y
	specular := NewVec3(1.0, 1.0, 1.0).MultiplyScalar(sp)
	re := reflectColor.MultiplyScalar(mat.albedo.z)
	return lit.Add(specular).Add(re)
}

// Reflect : Compute the illumination of a point using https://en.wikipedia.org/wiki/Phong_reflection_model
func Reflect(I Vec3, N Vec3) Vec3 {
	return I.Subtract(N.MultiplyScalar((I.DotProduct(N)) * 2.0))
}

// RayIntersect : http://www.lighthouse3d.com/tutorials/maths/ray-sphere-intersection/ (Check if a given ray intersects with a sphere)
func (s Sphere) RayIntersect(origin, direction Vec3, t float64) (float64, bool) {

	// Get the vector from the origin of the ray to the center of the sphere
	L := s.center.Subtract(origin)

	// Check how much the ray is in the direction of the sphere (https://math.stackexchange.com/questions/805954/what-does-the-dot-product-of-two-vectors-represent)
	tc := L.DotProduct(direction)

	// Sanity check, if its less than 0, its in the complete opposite direction (sphere is behind the ray origin)
	if tc < 0.0 {
		return t, false
	}

	// vec dot vec == sqrt(|vec|) (TODO: Double check understanding of this math)
	// [https://physics.info/vector-multiplication/#:~:text=Since%20the%20projection%20of%20a,square%20of%20that%20vector's%20magnitude.&text=Applying%20this%20corollary%20to%20the,vector%20with%20itself%20is%20one.]
	d2 := L.DotProduct(L) - (tc * tc)
	r2 := s.radius * s.radius

	// No intersection
	if d2 > r2 {
		return t, false
	}

	// Solve for magnitude
	tlc := math.Sqrt(r2 - d2)

	// Get interesection points
	t = tc - tlc
	t1 := tc + tlc

	if t < 0.0 {
		t = t1
	}

	if t < 0.0 {
		return t, false
	}

	return t, true
}
