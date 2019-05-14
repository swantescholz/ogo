package glmath

import (
	. "ogo/glmath/vec2"
	. "ogo/glmath/vec3"
	. "ogo/common"
)



func InterpolateBilinearVec2(A,B,C,D *Vec2, u,v Double) *Vec2 {
	p := A.Interpolate(B, u)
	q := C.Interpolate(D, u)
	return p.Interpolate(q, v)
}
func InterpolateBilinearVec3(A,B,C,D *Vec3, u,v Double) *Vec3 {
	p := A.Interpolate(B, u)
	q := C.Interpolate(D, u)
	return p.Interpolate(q, v)
}

func InterpolateHermiteVec3(A,U,B,V *Vec3, t Double) *Vec3 {
	a := A.Muls(2.0).Sub(B.Muls(2.0)).Add(U).Add(V)
	b := B.Muls(3.0).Sub(A.Muls(3.0)).Sub(U.Muls(2.0)).Sub(V)
	return a.Muls(t*t*t).Add(b.Muls(t*t)).Add(U.Muls(t)).Add(A)
}

func InterpolateHermiteD1Vec3(A,U,B,V *Vec3, t Double) *Vec3 {
	a := A.Muls(6.0).Sub(B.Muls(6.0)).Add(U.Muls(3.0)).Add(V.Muls(3.0))
	b := B.Muls(6.0).Sub(A.Muls(6.0)).Sub(U.Muls(4.0)).Sub(V.Muls(2.0))
	return a.Muls(t*t).Add(b.Muls(t)).Add(U)
}



