package color

import "math"
import "fmt"

type Color struct {
	R,G,B,A float64
}

var (
	White  = Col(1,1,1,1)
	Grey   = Col(.5,.5,.5,1)
	Black  = Col(0,0,0,1)
	Red    = Col(1,0,0,1)
	Green  = Col(0,1,0,1)
	Blue   = Col(0,0,1,1)
	Yellow = Col(1,1,0,1)
	Purple = Col(1,0,1,1)
	Cyan   = Col(0,1,1,1)
)

func Col(r,g,b,a float64) *Color {
	return &Color{r,g,b,a}
}

func (c *Color) Ptr() *float64 {
	return (*float64)(&c.R)
}
func (c *Color) Slc() []float32 {
	return []float32{float32(c.R),float32(c.G),float32(c.B),float32(c.A)}
}

func (c *Color) Copy() *Color {
	return Col(c.R, c.G, c.B, c.A)
}

func (c *Color) Equals(o *Color) bool {
	return c.R == o.R && c.G == o.G && c.B == o.B && c.A == o.A
}

func (c *Color) Set(o *Color) {
	c.R, c.G, c.B, c.A = o.R, o.G, o.B, o.A
}

func (c *Color) Length() float64 {
	return math.Sqrt((c.R * c.R) + (c.G * c.G) + (c.B * c.B) + (c.A * c.A))
}

func (c *Color) Length2() float64 {
	return (c.R * c.R) + (c.G * c.G) + (c.B * c.B) + (c.A * c.A)
}

func (c *Color) Distance(o *Color) float64 {
	return c.Sub(o).Length()
}

func (c *Color) Distance2(o *Color) float64 {
	return c.Sub(o).Length2()
}

func (c *Color) Unit() *Color {
	l := c.Length()
	return c.Divs(l)
}

func (c *Color) Add(o *Color) *Color {
	return Col(
		c.R+o.R,
		c.G+o.G,
		c.B+o.B,
		c.A+o.A)
}

func (c *Color) Adds(o float64) *Color {
	return Col(
		c.R+o,
		c.G+o,
		c.B+o,
		c.A+o)
}

func (c *Color) Sub(o *Color) *Color {
	return Col(
		c.R-o.R,
		c.G-o.G,
		c.B-o.B,
		c.A-o.A)
}

func (c *Color) Subs(o float64) *Color {
	return Col(
		c.R-o,
		c.G-o,
		c.B-o,
		c.A-o)
}

func (c *Color) Mul(o *Color) *Color {
	return Col(
		c.R*o.R,
		c.G*o.G,
		c.B*o.B,
		c.A*c.A)
}

func (c *Color) Muls(o float64) *Color {
	return Col(
		c.R*o,
		c.G*o,
		c.B*o,
		c.A*o)
}

func (c *Color) Div(o *Color) *Color {
	return Col(
		c.R/o.R,
		c.G/o.G,
		c.B/o.B,
		c.A/o.A)
}

func (c *Color) Divs(o float64) *Color {
	return c.Muls(float64(1)/o)
}

func (c *Color) String() string {
	return fmt.Sprintf("%v, %v, %v, %v", c.R,c.G,c.B,c.A)
}

