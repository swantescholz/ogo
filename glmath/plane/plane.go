package plane

import (
	. "ogo/glmath/vec3"
	. "ogo/common"
	"fmt"
)

//Plane with equation: N dot X + D = 0
type Plane struct {
	N *Vec3
	D Double
}

var (
	X   = New(V3(1,0,0),0)
	Y   = New(V3(0,1,0),0)
	Z   = New(V3(0,0,1),0)
	XY  = New(V3(1,1,0),0)
	XZ  = New(V3(1,0,1),0)
	YZ  = New(V3(0,1,1),0)
	XYZ = New(V3(1,1,1),0)
)

func New(normal *Vec3, distance Double) *Plane {
	return &Plane{N: normal.Copy(), D: distance}
}

func (t *Plane) String() string {
	return fmt.Sprintf("Plane(%v, %v)", t.N, t.D)
}

func ByPointAndNormal(p,n *Vec3) *Plane {
	return New(n, -1.0 * p.Dot(n))
}
func ByPoints(a,b,c *Vec3) *Plane {
	var n = c.Sub(b).Cross(a.Sub(b))
	var d Double = -1.0 * a.Dot(n)
	return &Plane{N: n, D: d}
}

func (t *Plane) Set(o *Plane) {
	t.N.Set(o.N)
	t.D = o.D
}
func (t *Plane) Copy() *Plane {
	return New(t.N, t.D)
}
//returns normalized copy
func (t *Plane) Unit() *Plane {
	finv := 1.0/t.N.Length()
	return New(t.N.Muls(finv), t.D*finv)
}
func (t *Plane) DotNormal(v *Vec3) Double {
	return t.N.Dot(v)
}
func (t *Plane) DotCoords(v *Vec3) Double {
	return t.N.Dot(v) + t.D
}
func (t *Plane) Distance(v *Vec3) Double {
	return t.Unit().DotCoords(v)
}
func (t *Plane) NearestPoint(v *Vec3) *Vec3 {
	return v.Sub(t.N.Muls(t.DotCoords(v)))
}


