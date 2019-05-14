package Color

import (
	"github.com/banthar/gl"
	"fmt"
	. "ogo/common"
	. "ogo/globals"
	"ogo/glmath/color24"
	"ogo/glmath/color32"
)


type Color struct {
	R,G,B Double
}

var (
	White  = &Color{1,1,1}
	Grey   = &Color{.5,.5,.5}
	Black  = &Color{0,0,0}
	Red    = &Color{1,0,0}
	Green  = &Color{0,1,0}
	Blue   = &Color{0,0,1}
	Yellow = &Color{1,1,0}
	Purple = &Color{1,0,1}
	Cyan   = &Color{0,1,1}
	Orange = &Color{1,0.5,0}
	Pink   = &Color{1,0.2,0.8}
	SkyBlue      = &Color{0.5,0.7,0.9}
	MellowOrange = &Color{0.9,0.7,0.5}
	ForestGreen  = &Color{0.247,0.498,0.373}
	Silver = &Color{0.784314,0.784314,0.784314}
	Gold   = &Color{0.862745,0.745098,0.0}
)


func Col0() *Color {
	return &Color{0,0,0}
}
func Col(r,g,b Double) *Color {
	return &Color{r,g,b}
}

func (c *Color) String() string {
	return fmt.Sprintf("%v, %v, %v", c.R,c.G,c.B)
}

func (c *Color) Gl() {gl.Color3dv(c.Slc())}
func (c *Color) Ptr() *float64 {return (*float64)(&c.R)}
func (c *Color) Slc32() []float32 {
	return []float32{float32(c.R),float32(c.G),float32(c.B)}
}
func (c *Color) Slc() []float64 {
	return []float64{float64(c.R),float64(c.G),float64(c.B)}
}
func (c *Color) Color24() *color24.Color24 {
	var rgba = c.Muls(256.0 - Epsilon)
	return &color24.Color24{uint8(rgba.R),uint8(rgba.G),uint8(rgba.B)}
}
func (c *Color) Color32() *color32.Color32 {
	var rgba = c.Muls(256.0 - Epsilon)
	return &color32.Color32{uint8(rgba.R),uint8(rgba.G),uint8(rgba.B),255}
}

func (c *Color) Copy() *Color {
	return Col(c.R, c.G, c.B)
}

func (c *Color) Equal(o Equaler) bool {
	return c.R == o.(*Color).R && c.G == o.(*Color).G && c.B == o.(*Color).B
}

func (c *Color) Set(o *Color) {
	c.R, c.G, c.B = o.R, o.G, o.B
}

func (c *Color) Interpolated(o *Color, t Double) *Color {
	return Col(
	(o.R - c.R) * t + c.R,
	(o.G - c.G) * t + c.G,
	(o.B - c.B) * t + c.B)
}
func (c *Color) Interpolate(o *Color, t Double) {
	c.R += (o.R - c.R) * t
	c.G += (o.G - c.G) * t
	c.B += (o.B - c.B) * t
}

func (c *Color) Length() Double {
	return c.Length2().Sqrt()
}

func (c *Color) Length2() Double {
	return (c.R * c.R) + (c.G * c.G) + (c.B * c.B)
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
	return Col(c.R+o.R, c.G+o.G, c.B+o.B)
}
func (c *Color) Addi(o *Color) {
	c.R += o.R
	c.G += o.G
	c.B += o.B
}

func (c *Color) Sub(o *Color) *Color {
	return Col(c.R-o.R, c.G-o.G, c.B-o.B)
}
func (c *Color) Subi(o *Color) {
	c.R += o.R
	c.G += o.G
	c.B += o.B
}

func (c *Color) Mul(o *Color) *Color {
	return Col( c.R*o.R, c.G*o.G, c.B*o.B)
}
func (c *Color) Muls(o Double) *Color {
	return Col( c.R*o, c.G*o, c.B*o)
}
func (c *Color) Muli(o *Color) {
	c.R *= o.R
	c.G *= o.G
	c.B *= o.B
}
func (c *Color) Mulsi(o Double) {
	c.R *= o
	c.G *= o
	c.B *= o
}

func (c *Color) Div(o *Color) *Color {
	return Col(c.R/o.R, c.G/o.G, c.B/o.B)
}
func (c *Color) Divs(o Double) *Color {
	return c.Muls(Double(1)/o)
}
func (c *Color) Divi(o *Color) {
	c.R /= o.R
	c.G /= o.G
	c.B /= o.B
}
func (c *Color) Divsi(o Double) {
	c.Mulsi(Double(1)/o)
}


