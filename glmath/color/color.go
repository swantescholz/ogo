package color

import (
	"github.com/banthar/gl"
	"fmt"
	. "ogo/common"
	. "ogo/globals"
	"ogo/glmath/color24"
	"ogo/glmath/color32"
)


type Color struct {
	R,G,B,A Double
}

var (
	Black0 = &Color{0,0,0,0}
	White  = &Color{1,1,1,1}
	Grey   = &Color{.5,.5,.5,1}
	Black  = &Color{0,0,0,1}
	Red    = &Color{1,0,0,1}
	Green  = &Color{0,1,0,1}
	Blue   = &Color{0,0,1,1}
	Yellow = &Color{1,1,0,1}
	Purple = &Color{1,0,1,1}
	Cyan   = &Color{0,1,1,1}
	Orange = &Color{1,0.5,0,1}
	Pink   = &Color{1,0.2,0.8,1}
	SkyBlue      = &Color{0.5,0.7,0.9,1}
	MellowOrange = &Color{0.9,0.7,0.5,1}
	ForestGreen  = &Color{0.247,0.498,0.373,1}
	Silver = &Color{0.784314,0.784314,0.784314,1}
	Gold   = &Color{0.862745,0.745098,0.0     ,1}
)


func Col0() *Color {
	return &Color{0,0,0,0}
}
func Col(r,g,b,a Double) *Color {
	return &Color{r,g,b,a}
}

func (c *Color) String() string {
	return fmt.Sprintf("%v, %v, %v, %v", c.R,c.G,c.B,c.A)
}

func (c *Color) Gl() {gl.Color4dv(c.Slc())}
func (c *Color) Ptr() *float64 {return (*float64)(&c.R)}
func (c *Color) Slc32() []float32 {
	return []float32{float32(c.R),float32(c.G),float32(c.B),float32(c.A)}
}
func (c *Color) Slc() []float64 {
	return []float64{float64(c.R),float64(c.G),float64(c.B),float64(c.A)}
}
func (c *Color) Color24() *color24.Color24 {
	var rgba = c.Muls(256.0 - Epsilon)
	return &color24.Color24{uint8(rgba.R),uint8(rgba.G),uint8(rgba.B)}
}
func (c *Color) Color32() *color32.Color32 {
	var rgba = c.Muls(256.0 - Epsilon)
	return &color32.Color32{uint8(rgba.R),uint8(rgba.G),uint8(rgba.B),uint8(rgba.A)}
}

func (c *Color) Copy() *Color {
	return Col(c.R, c.G, c.B, c.A)
}

func (c *Color) Equal(o Equaler) bool {
	return c.R == o.(*Color).R && c.G == o.(*Color).G && c.B == o.(*Color).B && c.A == o.(*Color).A
}

func (c *Color) Set(o *Color) {
	c.R, c.G, c.B, c.A = o.R, o.G, o.B, o.A
}

func (c *Color) Interpolated(o *Color, t Double) *Color {
	return Col(
	(o.R - c.R) * t + c.R,
	(o.G - c.G) * t + c.G,
	(o.B - c.B) * t + c.B,
	(o.A - c.A) * t + c.A)
}
func (c *Color) Interpolate(o *Color, t Double) {
	c.R += (o.R - c.R) * t
	c.G += (o.G - c.G) * t
	c.B += (o.B - c.B) * t
	c.A += (o.A - c.A) * t
}

func (c *Color) Length() Double {
	return c.Length2().Sqrt()
}

func (c *Color) Length2() Double {
	return (c.R * c.R) + (c.G * c.G) + (c.B * c.B) + (c.A * c.A)
}

func (c *Color) Distance(o *Color) Double {
	return c.Sub(o).Length()
}

func (c *Color) Distance2(o *Color) Double {
	return c.Sub(o).Length2()
}

func (c *Color) Unit() *Color {
	l := c.Length()
	return c.Divs(l)
}

func (c *Color) Add(o *Color) *Color {
	return Col(c.R+o.R, c.G+o.G, c.B+o.B, c.A+o.A)
}
func (c *Color) Addi(o *Color) {
	c.R += o.R
	c.G += o.G
	c.B += o.B
	c.A += o.A
}

func (c *Color) Sub(o *Color) *Color {
	return Col(c.R-o.R, c.G-o.G, c.B-o.B, c.A-o.A)
}
func (c *Color) Subi(o *Color) {
	c.R += o.R
	c.G += o.G
	c.B += o.B
	c.A += o.A
}

func (c *Color) Mul(o *Color) *Color {
	return Col( c.R*o.R, c.G*o.G, c.B*o.B, c.A*c.A)
}
func (c *Color) Muls(o Double) *Color {
	return Col( c.R*o, c.G*o, c.B*o, c.A*o)
}
func (c *Color) Muli(o *Color) {
	c.R *= o.R
	c.G *= o.G
	c.B *= o.B
	c.A *= c.A
}
func (c *Color) Mulsi(o Double) {
	c.R *= o
	c.G *= o
	c.B *= o
	c.A *= o
}

func (c *Color) Div(o *Color) *Color {
	return Col(c.R/o.R, c.G/o.G, c.B/o.B, c.A/o.A)
}
func (c *Color) Divs(o Double) *Color {
	return c.Muls(Double(1)/o)
}
func (c *Color) Divi(o *Color) {
	c.R /= o.R
	c.G /= o.G
	c.B /= o.B
	c.A /= o.A
}
func (c *Color) Divsi(o Double) {
	c.Mulsi(Double(1)/o)
}


