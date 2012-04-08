package vec2

import "math"
import "fmt"

type Vec2 struct {
	X, Y float64
}

func V2(x, y float64) *Vec2 {
	v := new(Vec2)
	v.X = x
	v.Y = y
	return v
}

func (v *Vec2) Ptr() *float64 {
	return (*float64)(&v.X)
}
func (v *Vec2) Slc() []float64 {
	return []float64{v.X,v.Y}
}

// Returns a new vector equal to 'v'
func (v *Vec2) Copy() *Vec2 {
	return V2(v.X, v.Y)
}

// Returns true if 'v' and 'o' are equal
func (v *Vec2) Equals(o *Vec2) bool {
	return v.X == o.X && v.Y == o.Y
}

func (v *Vec2) Set(o *Vec2) {
	v.X, v.Y = o.X, o.Y
}

// Returns the magnitude of 'v'
func (v *Vec2) Length() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y))
}

// Returns the squared magnitude of 'v'
func (v *Vec2) Length2() float64 {
	return (v.X * v.X) + (v.Y * v.Y)
}

// Returns the distance from 'v' to 'o'
func (v *Vec2) Distance(o *Vec2) float64 {
	return v.Sub(o).Length()
}

// Returns the squared distance from 'v' to 'o'
func (v *Vec2) Distance2(o *Vec2) float64 {
	return v.Sub(o).Length2()
}

// Returns the Dot product of 'v' and 'o'
func (v *Vec2) Dot(o *Vec2) float64 {
	return (v.X * o.X) + (v.Y * o.Y)
}

// Returns a new vector equal to 'v' but with a magnitude of 'l'
func (v *Vec2) SetLength(l float64) *Vec2 {
	u := v.Unit()
	return u.Muls(l)
}

// Returns a new vector equal to 'v' normalized
func (v *Vec2) Unit() *Vec2 {
	l := v.Length()
	return v.Divs(l)
}

func (v *Vec2) Add(o *Vec2) *Vec2 {
	return V2(
		v.X+o.X,
		v.Y+o.Y)
}

func (v *Vec2) Adds(o float64) *Vec2 {
	return V2(
		v.X+o,
		v.Y+o)
}

func (v *Vec2) Sub(o *Vec2) *Vec2 {
	return V2(
		v.X-o.X,
		v.Y-o.Y)
}

func (v *Vec2) Subs(o float64) *Vec2 {
	return V2(
		v.X-o,
		v.Y-o)
}

func (v *Vec2) Mul(o *Vec2) *Vec2 {
	return V2(
		v.X*o.X,
		v.Y*o.Y)
}

func (v *Vec2) Muls(o float64) *Vec2 {
	return V2(
		v.X*o,
		v.Y*o)
}

func (v *Vec2) Div(o *Vec2) *Vec2 {
	return V2(
		v.X/o.X,
		v.Y/o.Y)
}

func (v *Vec2) Divs(o float64) *Vec2 {
	return v.Muls(float64(1)/o)
}

func (v *Vec2) String() string {
	return fmt.Sprintf("%v, %v", v.X,v.Y)
}

