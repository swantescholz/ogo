package program

import (
	"github.com/banthar/gl"
	"fmt"
	"ogo/globals"
	"ogo/ass"
	"ogo/util"
)


var (
	GUsedProgram *Program
)

func init() {
	GUsedProgram = nil
}

type Program struct {
	id gl.Program
	VertShader, FragShader *Shader
	vparam, fparam string
}

func New(vert, frag string) *Program {
	ass.True(util.GlSupportsVersion(2,0))
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
	if GUsedProgram == nil {
		GUsedProgram = prog
	}
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

//call this AFTER setting view+projection
//call this BEFORE any other (in)direct uniform settings!
func (p *Program) Use () {
	if !globals.UseShader {
		panic("program used, although constants.UseShader==false")
	}
	if GCamera == nil {panic("no camera set for program")}
	GViewMatrix = GCamera.Matrix()
	GViewProjectionMatrix.Set(GViewMatrix.Mul(GProjectionMatrix))
	p.id.Use()
	p.SetViewMatrix(GViewMatrix)
	p.SetProjectionMatrix(GProjectionMatrix)
	p.SetViewProjectionMatrix(GViewProjectionMatrix)
	p.SetCamera(GCamera)
	p.ResetUniforms()
	GUsedProgram = p //remember me for global functions
}
func Unuse() {gl.ProgramUnuse()}
func (p *Program) String() string {
	return fmt.Sprintf("%v, %v, %v", p.id, p.VertShader.String(), p.FragShader.String())
}
func (p *Program) InfoLog() string {
	return fmt.Sprintf("%v,\n%v,\n%v", p.id.GetInfoLog(), p.VertShader.InfoLog(), p.FragShader.InfoLog())
}

func (p *Program) Failed() bool {
	a,b,c := p.id.GetInfoLog(), p.VertShader.InfoLog(), p.FragShader.InfoLog()
	if len(a) > 0 || len(b) > 0 || len(c) > 0 {
		return true
	}
	return false
}





