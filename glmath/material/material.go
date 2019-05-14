package material

import (
	"github.com/banthar/gl"
	"fmt"
	"ogo/glmath/color"
	"ogo/globals"
	"ogo/program"
	. "ogo/common"
)


type Material struct {
	A *color.Color //ambient
	B *color.Color //diffuss
	C *color.Color //specular
	S Double      //shininess
}

//satisfy program.IMaterial:
func (m *Material) Emission () *color.Color  {return &color.Color{0,0,0,0}}
func (m *Material) Ambient  () *color.Color  {return m.A}
func (m *Material) Diffuse  () *color.Color  {return m.B}
func (m *Material) Specular () *color.Color  {return m.C}
func (m *Material) Shininess() Double       {return m.S*128.0}

var (
	Zero  = New(color.Black0 , color.Black0 , color.Black0 , 0.0)
	White = New(color.White, color.White, color.White, 0.5)
	Grey  = New(color.Grey , color.Grey , color.White, 0.5)
	Black = New(color.Black, color.Black, color.White, 0.5)
	Red   = New(color.Red  , color.Red  , color.White, 0.5)
	Green = New(color.Green, color.Green, color.White, 0.5)
	Blue  = New(color.Blue , color.Blue , color.White, 0.5)
	Yellow = New(color.Yellow , color.Yellow , color.White, 0.5)
	Purple = New(color.Purple , color.Purple , color.White, 0.5)
	Cyan   = New(color.Cyan   , color.Cyan   , color.White, 0.5)
	Orange = New(color.Orange , color.Orange , color.White, 0.5)
	Pink   = New(color.Pink   , color.Pink   , color.White, 0.5)
	Sky    = New(color.SkyBlue      , color.SkyBlue      , color.White, 0.5)
	Mellow = New(color.MellowOrange , color.MellowOrange , color.White, 0.5)
	Forest = New(color.ForestGreen  , color.ForestGreen  , color.White, 0.5)
	Silver = New(color.Silver , color.Silver , color.White, 0.9)
	Gold   = New(color.Gold   , color.Gold   , color.White, 0.8)
)

func New(a,b,c *color.Color, s Double) *Material {
	m := new(Material)
	m.A,m.B,m.C,m.S = a,b,c,s
	return m
}

func (m *Material) String() string {
	return fmt.Sprintf("%v, %v, %v, %v", m.A,m.B,m.C,m.S)
}

func (m *Material) Use() {
	if globals.UseShader {
		program.SetFrontMaterial(m)
	} else {
		gl.Materialfv(gl.FRONT, gl.EMISSION, color.Black0.Slc32())
		gl.Materialfv(gl.FRONT, gl.AMBIENT  , m.A.Slc32())
		gl.Materialfv(gl.FRONT, gl.DIFFUSE  , m.B.Slc32())
		gl.Materialfv(gl.FRONT, gl.SPECULAR , m.C.Slc32())
		gl.Materialf (gl.FRONT, gl.SHININESS, float32(m.S*128.0))
	}
}

func (m *Material) Colors() (amb, diff, spec *color.Color) {
	amb, diff, spec = m.A, m.B, m.C
	return
}
func (m *Material) Copy() *Material {
	return New(m.A.Copy(), m.B.Copy(), m.C.Copy(), m.S)
}
func (m *Material) Equals(o *Material) bool {
	return m.A == o.A && m.B == o.B && m.C == o.C && m.S == o.S
}
func (m *Material) Set(o *Material) {
	m.A, m.B, m.C, m.S = o.A, o.B, o.C, o.S
}

func (m *Material) Add(o *Material) *Material {
	return New(m.A.Add(o.A), m.B.Add(o.B), m.C.Add(o.C), m.S+o.S)
}
func (m *Material) Addi(o *Material) {
	m.A.Addi(o.A)
	m.B.Addi(o.B)
	m.C.Addi(o.C)
	m.S += o.S
}

func (m *Material) Sub(o *Material) *Material {
	return New( m.A.Sub(o.A), m.B.Sub(o.B), m.C.Sub(o.C), m.S-o.S)
}
func (m *Material) Subi(o *Material) {
	m.A.Subi(o.A)
	m.B.Subi(o.B)
	m.C.Subi(o.C)
	m.S -= o.S
}

func (m *Material) Muls(o Double) *Material {
	return New( m.A.Muls(o), m.B.Muls(o), m.C.Muls(o), m.S*o)
}
func (m *Material) Mulsi(o Double) {
	m.A.Mulsi(o)
	m.B.Mulsi(o)
	m.C.Mulsi(o)
	m.S *= o
}

func (m *Material) Divs(o Double) *Material {
	return m.Muls(Double(1)/o)
}
func (m *Material) Divsi(o Double) {
	m.Mulsi(Double(1)/o)
}

