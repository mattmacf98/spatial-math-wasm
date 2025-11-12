// This serves as the dummy version of what spatial math could look like if we take out the dependency on the rdk utils (which is not WASM compilable)
package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
	"gonum.org/v1/gonum/num/quat"
)

const orientationVectorPoleRadius = 0.0001
const defaultAngleEpsilon = 1e-4

type R4AA struct {
	Theta float64 `json:"th"`
	RX    float64 `json:"x"`
	RY    float64 `json:"y"`
	RZ    float64 `json:"z"`
}

func (r4 *R4AA) Normalize() {
	norm := math.Sqrt(r4.RX*r4.RX + r4.RY*r4.RY + r4.RZ*r4.RZ)
	if norm == 0.0 { // prevent division by 0
		panic("cannot normalize R4AA, divide by zero")
	}
	r4.RX /= norm
	r4.RY /= norm
	r4.RZ /= norm
}

func (r4 *R4AA) ToQuat() quat.Number {
	sinA := math.Sin(r4.Theta / 2)
	// Ensure that point xyz is on the unit sphere
	r4.Normalize()

	// Get the unit-sphere components
	ax := r4.RX * sinA
	ay := r4.RY * sinA
	az := r4.RZ * sinA
	w := math.Cos(r4.Theta / 2)
	return quat.Number{w, ax, ay, az}
}

type OrientationVector struct {
	Theta float64 `json:"th"`
	OX    float64 `json:"x"`
	OY    float64 `json:"y"`
	OZ    float64 `json:"z"`
}

// QuatToOV converts a quaternion to an orientation vector.
func QuatToOV(q quat.Number) *OrientationVector {
	xAxis := quat.Number{0, -1, 0, 0}
	zAxis := quat.Number{0, 0, 0, 1}
	ov := &OrientationVector{}
	// Get the transform of our +X and +Z points
	newX := quat.Mul(quat.Mul(q, xAxis), quat.Conj(q))
	newZ := quat.Mul(quat.Mul(q, zAxis), quat.Conj(q))
	ov.OX = newZ.Imag
	ov.OY = newZ.Jmag
	ov.OZ = newZ.Kmag

	// The contents of ov.newX.Kmag are not in radians but we can use angleEpsilon anyway to check how close we are to
	// the pole because it's a convenient small number
	if 1-math.Abs(newZ.Kmag) > orientationVectorPoleRadius {
		v1 := mgl64.Vec3{newZ.Imag, newZ.Jmag, newZ.Kmag}
		v2 := mgl64.Vec3{newX.Imag, newX.Jmag, newX.Kmag}

		// Get the vector normal to the local-x, local-z, origin plane
		norm1 := v1.Cross(v2)

		// Get the vector normal to the global-z, local-z, origin plane
		norm2 := v1.Cross(mgl64.Vec3{zAxis.Imag, zAxis.Jmag, zAxis.Kmag})

		// For theta, we find the angle between the planes defined by local-x, global-z, origin and local-x, local-z, origin
		cosTheta := norm1.Dot(norm2) / (norm1.Len() * norm2.Len())
		// Account for floating point error
		if cosTheta > 1 {
			cosTheta = 1
		}
		if cosTheta < -1 {
			cosTheta = -1
		}

		theta := math.Acos(cosTheta)
		if theta > orientationVectorPoleRadius {
			// Acos will always produce a positive number, we need to determine directionality of the angle
			// We rotate newZ by -theta around the newX axis and see if we wind up coplanar with local-x, global-z, origin
			// If so theta is negative, otherwise positive
			// An R4AA is a convenient way to rotate a point by an amount around an arbitrary axis
			aa := R4AA{-theta, ov.OX, ov.OY, ov.OZ}
			q2 := aa.ToQuat()
			testZ := quat.Mul(quat.Mul(q2, zAxis), quat.Conj(q2))
			norm3 := v1.Cross(mgl64.Vec3{testZ.Imag, testZ.Jmag, testZ.Kmag})
			cosTest := norm1.Dot(norm3) / (norm1.Len() * norm3.Len())
			if 1-cosTest < defaultAngleEpsilon*defaultAngleEpsilon {
				ov.Theta = -theta
			} else {
				ov.Theta = theta
			}
		} else {
			ov.Theta = 0
		}
	} else {
		// Special case for when we point directly along the Z axis
		// Get the vector normal to the local-x, global-z, origin plane
		ov.Theta = -math.Atan2(newX.Jmag, -newX.Imag)
		if newZ.Kmag < 0 {
			ov.Theta = -math.Atan2(newX.Jmag, newX.Imag)
		}
	}
	// the IEEE 754 Standard for Floating-Points allows both negative and positive zero representations.
	// If one of the above conditions casts ov.Theta to -0, transform it to +0 for consistency.
	// String comparisons of floating point numbers may fail otherwise.
	if ov.Theta == -0. {
		ov.Theta = 0.
	}

	return ov
}
