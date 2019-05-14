package vec3

import (
	"github.com/banthar/gl"
	. "ogo/common"
	"math"
	"fmt"
)

type Vec3 struct {
	X, Y, Z Double
}

func V3(x, y, z Double) *Vec3 {
	return &Vec3{x,y,z}
}

func (v *Vec3) String() string {
	return fmt.Sprintf("Vec3(%v, %v, %v)", v.X,v.Y,v.Z)
}

func (v *Vec3) Gl() {
	gl.Vertex3dv(v.Slc())
}

func (v *Vec3) Ptr() *float64 {
	return (*float64)(&v.X)
}
func (v *Vec3) Slc32() []float32 {
	return []float32{float32(v.X),float32(v.Y),float32(v.Z)}
}
func (v *Vec3) Slc() []float64 {
	return []float64{float64(v.X),float64(v.Y),float64(v.Z)}
}

func (v *Vec3) Copy() *Vec3 {
	return V3(v.X, v.Y, v.Z)
}

func (v *Vec3) Equal(o Equaler) bool {
	return v.X == o.(*Vec3).X && v.Y == o.(*Vec3).Y && v.Z == o.(*Vec3).Z
}

func (v *Vec3) Set(o *Vec3) {
	v.X, v.Y, v.Z = o.X, o.Y, o.Z
}

func (v *Vec3) Interpolate(o *Vec3, t Double) *Vec3 {
	return v.Add(o.Sub(v).Muls(t))
}

func (v *Vec3) Length() Double {
	return Double(math.Sqrt(float64(v.Length2())))
}

func (v *Vec3) Length2() Double {
	return (v.X * v.X) + (v.Y * v.Y) + (v.Z * v.Z)
}

func (v *Vec3) Distance(o *Vec3) Double {
	return v.Sub(o).Length()
}

func (v *Vec3) Distance2(o *Vec3) Double {
	return v.Sub(o).Length2()
}

func (v *Vec3) Dot(o *Vec3) Double {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

func (v *Vec3) Cross(o *Vec3) *Vec3 {
	return V3(
		(v.Y*o.Z)-(v.Z*o.Y),
		(v.Z*o.X)-(v.X*o.Z),
		(v.X*o.Y)-(v.Y*o.X))
}

// Returns a vector reflected by 'norm'
func (v *Vec3) Reflect(norm *Vec3) *Vec3 {
	distance := Double(2) * v.Dot(norm)
	return V3(v.X-distance*norm.X,
		v.Y-distance*norm.Y,
		v.Z-distance*norm.Z)
}

// Returns a normalized vector of 'o' perpendicularized to 'v'
func (v *Vec3) Perp(o *Vec3) *Vec3 {
	perp := v.Cross(o)
	if perp.Length2() == 0 {
		perp = v.Cross(V3(0, 1, 0))
		if perp.Length2() == 0 {
			perp = v.Cross(V3(1, 0, 0))
		}
	}
	perp = perp.Cross(v)
	return perp.Normalized()
}

//returns angle between v and o as radian
func (v *Vec3) Angle(o *Vec3) Double {
	return (v.Dot(o) / (v.Length2() * o.Length2()).Sqrt()).Acos()
}

func (v *Vec3) Unit() *Vec3 {
	return v.Normalized()
}
func (v *Vec3) Normalized() *Vec3 {
	return v.Divs(v.Length())
}
func (v *Vec3) Normalize() {
	v.Divsi(v.Length())
}

func (v *Vec3) Add(o *Vec3) *Vec3 {
	return V3(v.X+o.X, v.Y+o.Y, v.Z+o.Z)
}
func (v *Vec3) Addi(o *Vec3) *Vec3 {
	v.X += o.X
	v.Y += o.Y
	v.Z += o.Z
	return v
}

func (v *Vec3) Sub(o *Vec3) *Vec3 {
	return V3(v.X-o.X, v.Y-o.Y, v.Z-o.Z)
}
func (v *Vec3) Subi(o *Vec3) *Vec3 {
	v.X -= o.X
	v.Y -= o.Y
	v.Z -= o.Z
	return v
}

func (v *Vec3) Mul(o *Vec3) *Vec3 {
	return V3(v.X*o.X, v.Y*o.Y, v.Z*o.Z)
}
func (v *Vec3) Muls(o Double) *Vec3 {
	return V3(v.X*o, v.Y*o, v.Z*o)
}
func (v *Vec3) Muli(o *Vec3) *Vec3 {
	v.X *= o.X
	v.Y *= o.Y
	v.Z *= o.Z
	return v
}
func (v *Vec3) Mulsi(o Double) *Vec3 {
	v.X *= o
	v.Y *= o
	v.Z *= o
	return v
}

func (v *Vec3) Div(o *Vec3) *Vec3 {
	return V3(v.X/o.X, v.Y/o.Y, v.Z/o.Z)
}
func (v *Vec3) Divs(o Double) *Vec3 {
	return v.Muls(Double(1)/o)
}
func (v *Vec3) Divi(o *Vec3) *Vec3 {
	v.X /= o.X
	v.Y /= o.Y
	v.Z /= o.Z
	return v
}
func (v *Vec3) Divsi(o Double) *Vec3 {
	v.Mulsi(Double(1)/o)
	return v
}

func Vec3RandCube() *Vec3 {
	return V3(RandDouble(),RandDouble(),RandDouble())
}
func Vec3RandUnit() *Vec3 {
	v := Vec3RandCube()
	for v.Length2() > 1.0 {
		v = Vec3RandCube()
	}
	return v.Unit()
}


