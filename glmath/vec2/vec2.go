package vec2

import (
	"github.com/banthar/gl"
	"math"
	"fmt"
	. "ogo/common"
)

type Vec2 struct {
	X, Y Double
}

func V2(x, y Double) *Vec2 {
	return &Vec2{x,y}
}

func (v *Vec2) String() string {
	return fmt.Sprintf("Vec2(%v, %v)", v.X,v.Y)
}

func (v *Vec2) Gl() {gl.Vertex2dv(v.Slc())}
func (v *Vec2) GlTex() {gl.TexCoord2dv(v.Slc())}
func (v *Vec2) Ptr() *float64 {return (*float64)(&v.X)}
func (v *Vec2) Slc32() []float32 {return []float32{float32(v.X),float32(v.Y)}}
func (v *Vec2) Slc() []float64 {return []float64{float64(v.X),float64(v.Y)}}
func (v *Vec2) Complex128() complex128 {return complex(v.X, v.Y)}
func (v *Vec2) Copy() *Vec2 {return V2(v.X, v.Y)}

func (v *Vec2) Equal(o Equaler) bool {
	return v.X == o.(*Vec2).X && v.Y == o.(*Vec2).Y
}

func (v *Vec2) Set(o *Vec2) {
	v.X, v.Y = o.X, o.Y
}

func (v *Vec2) Interpolate(o *Vec2, t Double) *Vec2 {
	return v.Add(o.Sub(v).Muls(t))
}

//-----------


func (v *Vec2) Length() Double {
	return Double(math.Sqrt(float64((v.X * v.X) + (v.Y * v.Y))))
}

func (v *Vec2) Length2() Double {
	return (v.X * v.X) + (v.Y * v.Y)
}

func (v *Vec2) Distance(o *Vec2) Double {
	return v.Sub(o).Length()
}

func (v *Vec2) Distance2(o *Vec2) Double {
	return v.Sub(o).Length2()
}

func (v *Vec2) Dot(o *Vec2) Double {
	return (v.X * o.X) + (v.Y * o.Y)
}

func (v *Vec2) Unit() *Vec2 {
	return v.Normalized()
}
func (v *Vec2) Normalized() *Vec2 {
	return v.Divs(v.Length())
}

func (v *Vec2) Normalize() {
	v.Divsi(v.Length())
}

//--------------

func (v *Vec2) Add(o *Vec2) *Vec2 {
	return V2(v.X + o.X, v.Y + o.Y)
}
func (v *Vec2) Addi(o *Vec2) *Vec2 {
	v.X += o.X
	v.Y += o.Y
	return v
}

func (v *Vec2) Sub(o *Vec2) *Vec2 {
	return V2(v.X - o.X, v.Y - o.Y)
}
func (v *Vec2) Subi(o *Vec2) *Vec2 {
	v.X -= o.X
	v.Y -= o.Y
	return v
}

func (v *Vec2) Mul(o *Vec2) *Vec2 {
	return V2(v.X * o.X, v.Y * o.Y)
}
func (v *Vec2) Muli(o *Vec2) *Vec2 {
	v.X *= o.X
	v.Y *= o.Y
	return v
}

func (v *Vec2) Muls(o Double) *Vec2 {
	return V2(v.X * o, v.Y * o)
}
func (v *Vec2) Mulsi(o Double) *Vec2 {
	v.X *= o
	v.Y *= o
	return v
}

func (v *Vec2) Div(o *Vec2) *Vec2 {
	return V2(v.X / o.X, v.Y / o.Y)
}
func (v *Vec2) Divi(o *Vec2) *Vec2 {
	v.X /= o.X
	v.Y /= o.Y
	return v
}

func (v *Vec2) Divs(o Double) *Vec2 {
	return v.Muls(Double(1.0)/o)
}
func (v *Vec2) Divsi(o Double) *Vec2 {
	v.Mulsi(Double(1.0)/o)
	return v
}

func Vec2RandSquare() *Vec2 {
	return V2(RandDouble(), RandDouble())
}
func Vec2RandUnit() *Vec2 {
	v := Vec2RandSquare()
	for v.Length2() > 1.0 {
		v = Vec2RandSquare()
	}
	return v.Unit()
}

var (
	Zero = &Vec2{0,0}
	One = &Vec2{1,0}
	One2 = &Vec2{1,1}
)

