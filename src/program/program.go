package program

import (
	"github.com/banthar/gl"
	"fmt"
	"strings"
	"io/ioutil"
	"glmath/color"
	"glmath/vec2"
	"glmath/vec3"
	"glmath/mat3"
	"glmath/mat4"
)


var (
	GDefaultShaderPathStart     = "res/shader/"
	GDefaultVertShaderPathStart = ""
	GDefaultFragShaderPathStart = ""
	GDefaultVertShaderPathEnd   = ".vert"
	GDefaultFragShaderPathEnd   = ".frag"
)
func SetDefaultShaderSearchPaths(s0, sv, sf, ev, ef string) {
	GDefaultShaderPathStart     = s0
	GDefaultVertShaderPathStart = sv
	GDefaultFragShaderPathStart = sf
	GDefaultVertShaderPathEnd   = ev
	GDefaultFragShaderPathEnd   = ef
}

type Shader struct {
	id gl.Shader
	isVert bool
}
func NewShader(isvert bool) *Shader {
	type_ := gl.GLenum(gl.FRAGMENT_SHADER)
	if isvert {type_ = gl.VERTEX_SHADER}
	id := gl.CreateShader(type_)
	sh := &Shader{id, isvert}
	return sh
}
func (s *Shader) Destroy() {s.id.Delete()}
func (s *Shader) IsVertShader() bool {return  s.isVert}
func (s *Shader) IsFragShader() bool {return !s.isVert}
func (s *Shader) InfoLog() string {return s.id.GetInfoLog()}
func (s *Shader) GetSource () string {return s.id.GetSource ()}
func (s *Shader) Source    (source string) {s.id.Source(source)}
func (s *Shader) Compile() {s.id.Compile()}
func (s *Shader) String() string {
	return fmt.Sprintf("%v", s.id)
}

func stringEndsWith(s, end string) bool {
	if len(end) > len(s) {return false}
	return s[len(s)-len(end):] == end
}
func ExtendVertShaderSourcePath(srcpath string) string {
	return GDefaultShaderPathStart + GDefaultVertShaderPathStart +
		srcpath + GDefaultVertShaderPathEnd
}
func ExtendFragShaderSourcePath(srcpath string) string {
	return GDefaultShaderPathStart + GDefaultFragShaderPathStart +
		srcpath + GDefaultFragShaderPathEnd
}

func PreprocessShaderSource(src string) string {
	const sinclude = "//#include "
	for {
		pos := strings.Index(src, sinclude)
		if pos < 0 {break}
		pos += len(sinclude)
		inewline := strings.Index(src, "\n")
		var path = ""
		if inewline < 0 {
			path = src[pos:]
		} else {
			path = src[pos: inewline]
		}
		path = GDefaultShaderPathStart + strings.TrimSpace(path)
		pos -= len(sinclude)
		b,_ := ioutil.ReadFile(path)
		extrasrc := string(b)
		endsrc := ""
		if inewline >= 0 { endsrc = src[inewline:]}
		src = src[:pos] + extrasrc + endsrc
		println(src)
	}
	return src
}

func vertShaderSource(sh string) string {
	if !strings.Contains(sh, "\n") { // && !strings.Contains(sh, "main") {
		sh = ExtendVertShaderSourcePath(sh)
		b,e := ioutil.ReadFile(sh)
		if e != nil {
			panic(e)
		}
		return string(b)
	}
	return sh
}
func fragShaderSource(sh string) string {
	if !strings.Contains(sh, "\n") { // && !strings.Contains(sh, "main") {
		sh = ExtendFragShaderSourcePath(sh)
		b,e := ioutil.ReadFile(sh)
		if e != nil {
			panic(e)
		}
		return string(b)
	}
	return sh
}

type Program struct {
	id gl.Program
	VertShader, FragShader *Shader
	vparam, fparam string
}

func New(vert, frag string) *Program {
	vsh := NewShader(true)
	fsh := NewShader(false)
	vsrc, fsrc := vertShaderSource(vert), fragShaderSource(frag)
	vsrc = PreprocessShaderSource(vsrc)
	fsrc = PreprocessShaderSource(fsrc)
	vsh.Source (vsrc); fsh.Source (fsrc)
	vsh.Compile()    ; fsh.Compile()
	prog := new(Program)
	prog.vparam, prog.fparam = vert, frag //save params
	prog.id = gl.CreateProgram()
	prog.AttachShader(vsh)
	prog.AttachShader(fsh)
	prog.Link()
	return prog
}
func (p *Program) Reload() {
	vsrc, fsrc := vertShaderSource(p.vparam), fragShaderSource(p.fparam)
	vsrc = PreprocessShaderSource(vsrc)
	fsrc = PreprocessShaderSource(fsrc)
	p.id.DetachShader(p.VertShader.id)
	p.id.DetachShader(p.FragShader.id)
	p.VertShader.Source (vsrc); p.FragShader.Source (fsrc)
	p.VertShader.Compile()    ; p.FragShader.Compile()
	p.AttachShader(p.VertShader)
	p.AttachShader(p.FragShader)
	p.Link()
}
func (p *Program) Destroy() {
	p.id.Delete()
	p.VertShader.Destroy()
	p.FragShader.Destroy()
}
func (p *Program) AttachShader(sh *Shader) {
	if sh.IsVertShader() {
		p.VertShader = sh
	} else {p.FragShader = sh}
	p.id.AttachShader(sh.id)
}
func (p *Program) Link() {p.id.Link()}
func (p *Program) Use () {
	p.id.Use()
}
func Unuse() {gl.ProgramUnuse()}
func (p *Program) String() string {
	return fmt.Sprintf("%v, %v, %v", p.id, p.VertShader.String(), p.FragShader.String())
}
func (p *Program) InfoLog() string {
	return fmt.Sprintf("%v,\n%v,\n%v", p.id.GetInfoLog(), p.VertShader.InfoLog(), p.FragShader.InfoLog())
}

func (p *Program) GetUniformLocation(name string) gl.UniformLocation {
	return p.id.GetUniformLocation(name)
}
func (p *Program) UniformVec2(name string, v *vec2.Vec2) {
	loc := p.GetUniformLocation(name)
	loc.Uniform2f(float32(v.X), float32(v.Y))
}
func (p *Program) UniformVec3(name string, v *vec3.Vec3) {
	loc := p.GetUniformLocation(name)
	loc.Uniform3f(float32(v.X), float32(v.Y), float32(v.Z))
}
func (p *Program) UniformColor(name string, c *color.Color) {
	loc := p.GetUniformLocation(name)
	loc.Uniform4f(float32(c.R), float32(c.G), float32(c.B), float32(c.A))
}
func (p *Program) UniformColor3(name string, c *color.Color) {
	loc := p.GetUniformLocation(name)
	loc.Uniform3f(float32(c.R), float32(c.G), float32(c.B))
}
func (p *Program) UniformMat3(name string, m *mat3.Mat3) {
	loc := p.GetUniformLocation(name)
	loc.Uniform1fv(m.Fv32())
}
func (p *Program) UniformMat4(name string, m *mat4.Mat4) {
	loc := p.GetUniformLocation(name)
	//loc.Uniform1fv(m.Fv32())
	loc.UniformMatrix4fv(m.Fv32(),1,false)
	//f := m.Fv32()
	//for i := range f {
	//	fmt.Println(f[i])
	//}
}
func (p *Program) SetModelMatrix(m *mat4.Mat4) {
	//println("------------------------")
	p.UniformMat4("uModelMatrix",m) //non-standard!
	//println("------------------------")
}
func (p *Program) SetViewMatrix(m *mat4.Mat4) {
	p.UniformMat4("uViewMatrix",m) //non-standard!
}
func (p *Program) SetModelViewMatrix(m *mat4.Mat4) {
	p.UniformMat4("gl_ModelViewMatrix",m) 
}
func (p *Program) SetModelViewProjectionMatrix(m *mat4.Mat4) {
	p.UniformMat4("gl_ModelViewProjectionMatrix",m)
}
func (p *Program) SetProjectionMatrix(m *mat4.Mat4) {
	p.UniformMat4("gl_ProjectionMatrix",m)
}


