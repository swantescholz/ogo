package material

//import "math"
import "fmt"
import "glmath/color"
import "github.com/banthar/gl"

type Material struct {
	A *color.Color //ambient
	B *color.Color //diffuss
	C *color.Color //specular
	S float64      //shininess
}

var (
	White = New(color.White, color.White, color.White, 0.5)
	Grey  = New(color.Grey , color.Grey , color.White, 0.5)
	Black = New(color.Black, color.Black, color.White, 0.5)
	Red   = New(color.Red  , color.Red  , color.White, 0.5)
	Green = New(color.Green, color.Green, color.White, 0.5)
	Blue  = New(color.Blue , color.Blue , color.White, 0.5)
)

func New(a,b,c *color.Color, s float64) *Material {
	m := new(Material)
	m.A,m.B,m.C,m.S = a,b,c,s
	return m
}

func (m *Material) Use() {
	gl.Materialfv(gl.FRONT, gl.AMBIENT , m.A.Slc());
	gl.Materialfv(gl.FRONT, gl.DIFFUSE , m.B.Slc());
	gl.Materialfv(gl.FRONT, gl.SPECULAR, m.C.Slc());
	gl.Materialf (gl.FRONT, gl.SHININESS, float32(m.S*128.0));
}

// Returns a new vector equal to 'v'
func (m *Material) Copy() *Material {
	return New(m.A, m.B, m.C, m.S)
}

// Returns true if 'v' and 'o' are equal
func (m *Material) Equals(o *Material) bool {
	return m.A == o.A && m.B == o.B && m.C == o.C && m.S == o.S
}

func (m *Material) Set(o *Material) {
	m.A, m.B, m.C, m.S = o.A, o.B, o.C, o.S
}

func (m *Material) Add(o *Material) *Material {
	return New(
		m.A.Add(o.A),
		m.B.Add(o.B),
		m.C.Add(o.C),
		m.S+o.S)
}

func (m *Material) Adds(o float64) *Material {
	return New(
		m.A.Adds(o),
		m.B.Adds(o),
		m.C.Adds(o),
		m.S+o)
}

func (m *Material) Sub(o *Material) *Material {
	return New(
		m.A.Sub(o.A),
		m.B.Sub(o.B),
		m.C.Sub(o.C),
		m.S-o.S)
}

func (m *Material) Subs(o float64) *Material {
	return m.Adds(-o)
}


func (m *Material) Muls(o float64) *Material {
	return New(
		m.A.Muls(o),
		m.B.Muls(o),
		m.C.Muls(o),
		m.S*o)
}

func (m *Material) Divs(o float64) *Material {
	return m.Muls(float64(1)/o)
}

func (m *Material) String() string {
	return fmt.Sprintf("%v, %v, %v, %v", m.A,m.B,m.C,m.S)
}
