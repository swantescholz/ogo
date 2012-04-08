package vec3

import "math"
import "fmt"

type Vec3 struct {
	X, Y, Z float64
}

func V3(x, y, z float64) *Vec3 {
	v := new(Vec3)
	v.X = x
	v.Y = y
	v.Z = z
	return v
}

func (v *Vec3) Ptr() *float64 {
	return (*float64)(&v.X)
}
func (v *Vec3) Slc() []float64 {
	return []float64{v.X,v.Y,v.Z}
}

// Returns a new vector equal to 'v'
func (v *Vec3) Copy() *Vec3 {
	return V3(v.X, v.Y, v.Z)
}

// Returns true if 'v' and 'o' are equal
func (v *Vec3) Equals(o *Vec3) bool {
	return v.X == o.X && v.Y == o.Y && v.Z == o.Z
}

func (v *Vec3) Set(o *Vec3) {
	v.X, v.Y, v.Z = o.X, o.Y, o.Z
}

// Returns the magnitude of 'v'
func (v *Vec3) Length() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y) + (v.Z * v.Z))
}

// Returns the squared magnitude of 'v'
func (v *Vec3) Length2() float64 {
	return (v.X * v.X) + (v.Y * v.Y) + (v.Z * v.Z)
}

// Returns the distance from 'v' to 'o'
func (v *Vec3) Distance(o *Vec3) float64 {
	return v.Sub(o).Length()
}

// Returns the squared distance from 'v' to 'o'
func (v *Vec3) Distance2(o *Vec3) float64 {
	return v.Sub(o).Length2()
}

// Returns the Dot product of 'v' and 'o'
func (v *Vec3) Dot(o *Vec3) float64 {
	return (v.X * o.X) + (v.Y * o.Y) + (v.Z * o.Z)
}

// Returns the cross product of 'v' and 'o'
func (v *Vec3) Cross(o *Vec3) *Vec3 {
	return V3(
		(v.Y*o.Z)-(v.Z*o.Y),
		(v.Z*o.X)-(v.X*o.Z),
		(v.X*o.Y)-(v.Y*o.X))
}

// Returns a vector reflected by 'norm'
func (v *Vec3) Reflect(norm *Vec3) *Vec3 {
	distance := float64(2) * v.Dot(norm)
	return V3(v.X-distance*norm.X,
		v.Y-distance*norm.Y,
		v.Z-distance*norm.Z)
}

// Returns a normalized vector perpendicular to 'v'
/*func (v *Vec3) Perp() *Vec3 {
	perp := v.Cross(V3(-1, 0, 0))
	if perp.Length() == 0 {
		// If v is too close to -x try -y
		perp = v.Cross(V3(0, -1, 0))
	}
	return perp.Unit()
}//*/

// Returns a new vector equal to 'v' but with a magnitude of 'l'
func (v *Vec3) SetLength(l float64) *Vec3 {
	u := v.Unit()
	return u.Muls(l)
}

// Returns a new vector equal to 'v' normalized
func (v *Vec3) Unit() *Vec3 {
	l := v.Length()
	return v.Divs(l)
}

func (v *Vec3) Add(o *Vec3) *Vec3 {
	return V3(
		v.X+o.X,
		v.Y+o.Y,
		v.Z+o.Z)
}

func (v *Vec3) Adds(o float64) *Vec3 {
	return V3(
		v.X+o,
		v.Y+o,
		v.Z+o)
}

func (v *Vec3) Sub(o *Vec3) *Vec3 {
	return V3(
		v.X-o.X,
		v.Y-o.Y,
		v.Z-o.Z)
}

func (v *Vec3) Subs(o float64) *Vec3 {
	return V3(
		v.X-o,
		v.Y-o,
		v.Z-o)
}

func (v *Vec3) Mul(o *Vec3) *Vec3 {
	return V3(
		v.X*o.X,
		v.Y*o.Y,
		v.Z*o.Z)
}

func (v *Vec3) Muls(o float64) *Vec3 {
	return V3(
		v.X*o,
		v.Y*o,
		v.Z*o)
}

func (v *Vec3) Div(o *Vec3) *Vec3 {
	return V3(
		v.X/o.X,
		v.Y/o.Y,
		v.Z/o.Z)
}

func (v *Vec3) Divs(o float64) *Vec3 {
	return v.Muls(float64(1)/o)
}

func (v *Vec3) String() string {
	return fmt.Sprintf("%v, %v, %v", v.X,v.Y,v.Z)
}


