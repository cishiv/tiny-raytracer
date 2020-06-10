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
	albedo           Vec4
	specularExponent float64
	refractiveIndex  float64
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
func NewMaterial(diffuseColor Vec3, albedo Vec4, specularExponent, refractiveIndex float64) Material {
	return Material{
		diffuseColor,
		albedo,
		specularExponent,
		refractiveIndex,
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

	checkerboardDistance := float64(math.MaxFloat64)
	if math.Abs(direction.y) > 1E-3 {
		d := -(origin.y + 4) / direction.y
		pt := origin.Add(direction.MultiplyScalar(d))
		if d > 0.0 && math.Abs(pt.x) < 10.0 && pt.z < -10 && pt.z > -30.0 && d < sphereDistance {
			checkerboardDistance = d
			hit.Copy(pt)
			N.Copy(NewVec3(0, 1, 0))
			if (int(0.5*hit.x+1000)+int(0.5*hit.z))&1 == 1 {
				mat.diffuseColor.Copy(NewVec3(.3, .3, .3))
			} else {
				mat.diffuseColor.Copy(NewVec3(.3, .2, .1))
			}
			mat.refractiveIndex = 1
			mat.specularExponent = 50.0
			mat.albedo = NewVec4(1, 0.2, 0, 0)
		}
	}
	return math.Min(sphereDistance, checkerboardDistance) < 1000.0

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

	// offset
	// if the point drawn lines on the surface of the object, any ray will intersect the object (?)
	// prevents self occlusion when factored in
	p := 1E-3

	// Reflection
	// When intersecting the sphere, compute the reflected ray
	// And then just cast a ray in its direction
	// We don't need to normalize here, but just for safety...
	reflectDirection := Reflect(*direction, *N)
	var reflectOrigin Vec3
	if reflectDirection.DotProduct(*N) < 0 {
		reflectOrigin = point.Subtract(N.MultiplyScalar(p))
	} else {
		reflectOrigin = point.Add(N.MultiplyScalar(p))
	}

	reflectColor := CastRay(&reflectOrigin, &reflectDirection, spheres, lights, depth+1)

	// Refraction
	refractDirection := Refract(*direction, *N, mat.refractiveIndex, 1.0)
	var refractOrigin Vec3
	if refractDirection.DotProduct(*N) < 0 {
		refractOrigin = point.Subtract(N.MultiplyScalar(p))
	} else {
		refractOrigin = point.Add(N.MultiplyScalar(p))
	}

	refractColor := CastRay(&refractOrigin, &refractDirection, spheres, lights, depth+1)

	// Shadows + Diffuse
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

	// Quick math
	lit := mat.diffuseColor.MultiplyScalar(diffuseLightIntensity).MultiplyScalar(mat.albedo.x)
	sp := specularLightIntensity * mat.albedo.y
	specular := NewVec3(1.0, 1.0, 1.0).MultiplyScalar(sp)
	re := reflectColor.MultiplyScalar(mat.albedo.z)
	refract := refractColor.MultiplyScalar(mat.albedo.w)
	return lit.Add(specular).Add(re).Add(refract)
}

// Reflect : Compute the illumination of a point using https://en.wikipedia.org/wiki/Phong_reflection_model
func Reflect(I Vec3, N Vec3) Vec3 {
	return I.Subtract(N.MultiplyScalar((I.DotProduct(N)) * 2.0))
}

// Refract : Compute the refracted ray using Snells Law (https://en.wikipedia.org/wiki/Snell%27s_law)
func Refract(I Vec3, N Vec3, etaT float64, etaI float64) Vec3 {
	cosi := -math.Max(-1.0, math.Min(1.0, I.DotProduct(N)))

	if cosi < 0 {
		// If ray is inside the object, swap
		return Refract(I, N.Negate(), etaI, etaT)
	}

	eta := etaI / etaT

	k := 1 - eta*eta*(1-cosi*cosi)

	if k < 0 {
		return NewVec3(1.0, 0.0, 0.0)
	}

	return I.MultiplyScalar(eta).Add(N.MultiplyScalar((eta*cosi - math.Sqrt(k))))
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
