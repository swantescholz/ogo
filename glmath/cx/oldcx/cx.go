package oldcx

import (
	"math"
	. "fmt"
	. "ogo/common"
	"ogo/glmath/vec2"
)


func SetPrec(prec int) {}
func NewFloat(f Double) Double {return Double(f)}
func SetIterations(i int) {iterations = i}
const FloatType = "Double"

type Cx struct {
	X, Y Double
}

var a,b,cx,cy,x,y Double

var iterations int

func New(x, y Double) *Cx {
	return &Cx{x,y}
}
func NewFromString(x, y string) *Cx {
	Sscan(x, a)
	Sscan(y, b)
	return New(a,b)
}

func NewFromVec2(v *vec2.Vec2) *Cx {
	return New(v.X, v.Y)
}
func NewFromFloat(x,y Double) *Cx {
	return New(x,y)
}

func (v *Cx) Clear() {}

func (v *Cx) String() string {
	return Sprintf("Cx(%v, %v)", v.X,v.Y)
}

func (v *Cx) Complex128() complex128 {return complex(v.X, v.Y)}
func (v *Cx) Copy() *Cx {return New(v.X, v.Y)}

func (v *Cx) Equal(o Equaler) bool {
	return v.X == o.(*Cx).X && v.Y == o.(*Cx).Y
}

func (v *Cx) Set(o *Cx) {
	v.X, v.Y = o.X, o.Y
}
func (v *Cx) SetString(x,y string) {
	Sscan(x, &v.X)
	Sscan(y, &v.Y)
}
func (v *Cx) SetX(x Double) {v.X = x}
func (v *Cx) SetY(y Double) {v.Y = y}
func (v *Cx) SetPrec(prec int) {}

//-----------

func (v *Cx) Inverted() *Cx {
	var l2inv = 1.0 / v.Length2()
	return New(v.X * l2inv, -v.Y * l2inv)
}
func (v *Cx) Invert() {
	var l2inv = 1.0 / v.Length2()
	v.X *= l2inv
	v.Y *= -l2inv
}

func (v *Cx) Squared() *Cx {
	return New(v.X*v.X - v.Y*v.Y, 2.0*v.X*v.Y)
}
func (v *Cx) Square() {
	//*
	b = (v.X + v.Y)*(v.X - v.Y)
	v.Y = 2.0*v.X*v.Y
	v.X = b//*/
	//v.X, v.Y = v.X*v.X - v.Y*v.Y, 2.0*v.X*v.Y
}

func (v *Cx) DoMandel(c *Cx) int {
	x,y,cx,cy = v.X,v.Y, c.X,c.Y
	for i := 0;; i++ {
		a,b = x*x, y*y
		if a+b > 4.0 {
			return i
		}
		if i == iterations-1 {break}
		y = 2.0*x*y + cy
		x = a - b + cx
	}
	return iterations
}

func (v *Cx) Length() Double {
	return Double(math.Sqrt(float64((v.X * v.X) + (v.Y * v.Y))))
}

func (v *Cx) Length2() Double {
	return (v.X * v.X) + (v.Y * v.Y)
}
func (v *Cx) Length2Double() Double {
	return (v.X * v.X) + (v.Y * v.Y)
}

func (v *Cx) Distance(o *Cx) Double {
	return v.Sub(o).Length()
}

func (v *Cx) Distance2(o *Cx) Double {
	return v.Sub(o).Length2()
}

func (v *Cx) Dot(o *Cx) Double {
	return (v.X * o.X) + (v.Y * o.Y)
}

func (v *Cx) Normalized() *Cx {
	l := v.Length()
	return v.Divs(l)
}

func (v *Cx) Normalize() {
	var l = v.Length()
	v.Divsi(l)
}

//--------------

func (v *Cx) Add(o *Cx) *Cx {
	return New(v.X + o.X, v.Y + o.Y)
}
func (v *Cx) Addi(o *Cx) {
	v.X += o.X
	v.Y += o.Y
}
func (v *Cx) AddX(o Double) {v.X += o}
func (v *Cx) AddY(o Double) {v.Y += o}
func (v *Cx) MulX(o Double) {v.X *= o}
func (v *Cx) MulY(o Double) {v.Y *= o}

func (v *Cx) Sub(o *Cx) *Cx {
	return New(v.X - o.X, v.Y - o.Y)
}
func (v *Cx) Subi(o *Cx) {
	v.X -= o.X
	v.Y -= o.Y
}

func (v *Cx) Mul(o *Cx) *Cx {
	return New(v.X * o.X - v.Y * o.Y, v.X * o.Y + v.Y * o.X)
}
func (v *Cx) Muli(o *Cx) {
	oldx := v.X
	v.X = v.X * o.X - v.Y * o.Y
	v.Y = oldx * o.Y + v.Y * o.X
}

func (v *Cx) MuliComponents(o *Cx) {
	v.X *= o.X
	v.Y *= o.Y
}

func (v *Cx) MulsDouble(o Double) *Cx {return v.Muls(o)}
func (v *Cx) Muls(o Double) *Cx {
	return New(v.X * o, v.Y * o)
}
func (v *Cx) Mulsi(o Double) {
	v.X *= o
	v.Y *= o
}

func (v *Cx) Div(o *Cx) *Cx {
	var tmp = o.Inverted()
	tmp.Muli(v)
	return tmp
}
func (v *Cx) Divi(o *Cx) {
	v.Divsi(o.Length2())
	oldx := v.X
	v.X = v.X * o.X + v.Y * o.Y
	v.Y = oldx * -o.Y + v.Y * o.X
}

func (v *Cx) Divs(o Double) *Cx {
	return v.Muls(Double(1.0)/o)
}
func (v *Cx) Divsi(o Double) {
	v.Mulsi(Double(1.0)/o)
}
func (v *Cx) DivComponents(o *Cx) *Cx {
	return New(v.X/o.X, v.Y/o.Y)
}

var (
	Zero = &Cx{0,0}
	One = &Cx{1,0}
	One2 = &Cx{1,1}
)

