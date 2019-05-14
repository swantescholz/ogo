package cx

import (
	"github.com/salviati/gmp"
	. "fmt"
	. "ogo/common"
	"ogo/glmath/vec2"
)

func NewFloat(d Double) *gmp.Float {return gmp.NewFloat(d.F())}
func SetIterations(i int) {iterations = i}
const FloatType = "gmp.Float"

var (
	a,b,cx,cy,d,x,y *gmp.Float
	iterations int
)
func SetPrec(prec int) {
	a.SetPrec(prec)
	b.SetPrec(prec)
	cx.SetPrec(prec)
	cy.SetPrec(prec)
	d.SetPrec(prec)
	x.SetPrec(prec)
	y.SetPrec(prec)
}

func init() {
	a,b,x,y = NewFloat(0),NewFloat(0),NewFloat(0),NewFloat(0)
	cx,cy,d = NewFloat(0),NewFloat(0),NewFloat(0)
}

type Cx struct {
	X, Y *gmp.Float
}

func New(x, y Double) *Cx {
	var this = &Cx{NewFloat(x), NewFloat(y)}
	return this
}
func NewFromString(x, y string) *Cx {
	var this = New(0,0)
	this.SetString(x,y)
	return this
}
func NewFromFloat(x, y *gmp.Float) *Cx {
	var this = new(Cx)
	this.X = x.Copy()
	this.Y = y.Copy()
	return this
}

func NewFromVec2(v *vec2.Vec2) *Cx {
	return New(v.X, v.Y)
}

func (v *Cx) String() string {
	return Sprintf("Cx(%v, %v)", v.X,v.Y)
}

func (v *Cx) destroy() {
	v.X.Clear()
	v.Y.Clear()
}
func (v *Cx) Clear() {
	v.destroy()
}

func (v *Cx) Copy() *Cx {
	var n = new(Cx)
	n.X = v.X.Copy()
	n.Y = v.Y.Copy()
	return n
}

func (v *Cx) Equal(o Equaler) bool {
	return v.X.Cmp(o.(*Cx).X) == 0 && v.Y.Cmp(o.(*Cx).Y) == 0
}

func (v *Cx) Set(o *Cx) {
	v.X.Set(o.X)
	v.Y.Set(o.Y)
}
func (v *Cx) SetString(x, y string) {
	v.X.SetString(x)
	v.Y.SetString(y)
}

func (v *Cx) SetX(x *gmp.Float) {v.X.Set(x)}
func (v *Cx) SetY(y *gmp.Float) {v.Y.Set(y)}

func (v *Cx) SetPrec(prec int) {
	v.X.SetPrec(prec)
	v.Y.SetPrec(prec)
}
func (v *Cx) GetPrec() int {
	return v.X.GetPrec()
}

//-----------

func (v *Cx) Square() {
	a.Add(v.X,v.Y)
	b.Sub(v.X,v.Y)
	v.Y.Mul(v.X,v.Y)
	v.Y.Mul2Exp(v.Y,1)
	v.X.Mul(a,b)
}

func (v *Cx) DoMandel(c *Cx) int {
	x.Set(v.X); y.Set(v.Y)
	cx.Set(c.X); cy.Set(c.Y)
	return gmp.DoMandel(cx,cy,x,y,iterations)
	/*
	for i := 0;; i++ {
		a.Mul(x,x)
		b.Mul(y,y)
		if a.Double() + b.Double() > 4.0 {
			return i
		}
		if i == iterations-1 {break}
		y.Mul2Exp(y,1).Mul(y,x).Add(y,cy)
		x.Sub(a,b).Add(x,cx)
	}
	return iterations//*/
}

func (v *Cx) Length2Double() Double {
	a.Mul(v.X,v.X)
	b.Mul(v.Y,v.Y)
	a.Add(a,b)
	return Double(a.Double())
}

func (v *Cx) Addi(o *Cx) {
	v.X.Add(v.X,o.X)
	v.Y.Add(v.Y,o.Y)
}
func (v *Cx) AddX(o *gmp.Float) {v.X.Add(v.X,o)}
func (v *Cx) AddY(o *gmp.Float) {v.Y.Add(v.Y,o)}
func (v *Cx) MulX(o *gmp.Float) {v.X.Mul(v.X,o)}
func (v *Cx) MulY(o *gmp.Float) {v.Y.Mul(v.Y,o)}

func (v *Cx) Sub(o *Cx) *Cx {
	var a = v.Copy()
	a.X.Sub(a.X, o.X)
	a.Y.Sub(a.Y, o.Y)
	return a
}
func (v *Cx) Subi(o *Cx) {
	v.X.Sub(v.X, o.X)
	v.Y.Sub(v.Y, o.Y)
}

func (v *Cx) MuliComponents(o *Cx) {
	v.X.Mul(v.X, o.X)
	v.Y.Mul(v.Y, o.Y)
}
func (v *Cx) Muls(o *gmp.Float) *Cx {
	var a = v.Copy()
	a.X.Mul(a.X,o)
	a.Y.Mul(a.Y,o)
	return a
}
func (v *Cx) Mulsi(o *gmp.Float) {
	v.X.Mul(v.X,o)
	v.Y.Mul(v.Y,o)
}
func (v *Cx) MulsDouble(o Double) *Cx {
	var f = NewFloat(o)
	defer f.Clear()
	return v.Muls(f)
}

func (v *Cx) DivComponents(o *Cx) *Cx {
	var a = v.Copy()
	a.X.Div(a.X, o.X)
	a.Y.Div(a.Y, o.Y)
	return a
}

