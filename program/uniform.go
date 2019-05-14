package program

import (
	"github.com/banthar/gl"
	"ogo/glmath/color"
	"ogo/glmath/vec2"
	"ogo/glmath/vec3"
	"ogo/glmath/mat3"
	"ogo/glmath/mat4"
	"ogo/globals"
	"fmt"
	. "ogo/common"
)

var (
	GViewMatrix, GProjectionMatrix, GViewProjectionMatrix *mat4.Mat4
	GFarClipplane Double
	GCamera ICamera
)

func init() {
	GViewMatrix = mat4.New()
	GProjectionMatrix = mat4.New()
	GViewProjectionMatrix = mat4.New()
	GFarClipplane = 1024.0
	GCamera = nil
}

type IMaterial interface {
	Emission() *color.Color
	Ambient() *color.Color
	Diffuse() *color.Color
	Specular() *color.Color
	Shininess() Double
}
type ILight interface {
	Id() int
	Ambient() *color.Color
	Diffuse() *color.Color
	Specular() *color.Color
	Position() *vec3.Vec3
	HalfVector() *vec3.Vec3
	SpotDirection() *vec3.Vec3
	SpotExponent() Double
	SpotCutoff() Double
	SpotCosCutoff() Double
	ConstantAttenuation() Double
	LinearAttenuation() Double
	QuadraticAttenuation() Double
	Enabled() bool
	Positional() bool
}
type ICamera interface {
	Position() *vec3.Vec3
	Direction() *vec3.Vec3
	UpVector() *vec3.Vec3
	Matrix() *mat4.Mat4
}
type ITexture interface {
	Id() gl.Texture
}


//UNIFORMS -----------------------------------
func (p *Program) ResetUniforms() {
	for i := 0; i < globals.MaxLights; i++ {
		name := fmt.Sprintf("uLightSource[%v].enabled")
		p.UniformBool(name, false) //disable
	}
	p.SetTextureMatrix(mat4.Identity())
	p.SetModelMatrix(mat4.Identity())
	p.SetNormalMatrix(mat3.Identity())
	p.UniformFloat("uFarClipplane", GFarClipplane)
}

func (p *Program) GetUniformLocation(name string) gl.UniformLocation {
	
	i := p.id.GetUniformLocation(name)
	//println("uniloc", i)
	//println("ldd", name)
	return i
}

func (p *Program) UniformBool(name string, b bool) {
	i := 0
	if b {i = 1}
	loc := p.GetUniformLocation(name)
	loc.Uniform1i(i)
}
func (p *Program) UniformInt(name string, i int) {
	loc := p.GetUniformLocation(name)
	loc.Uniform1i(i)
}
func (p *Program) UniformFloat(name string, f Double) {
	loc := p.GetUniformLocation(name)
	loc.Uniform1f(float32(f))
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
func (p *Program) UniformMaterial(name string, m IMaterial) {
	p.UniformColor(name + ".emission" , m.Emission())
	p.UniformColor(name + ".ambient"  , m.Ambient())
	p.UniformColor(name + ".diffuse"  , m.Diffuse())
	p.UniformColor(name + ".specular" , m.Specular())
	p.UniformFloat(name + ".shininess", m.Shininess())
}
func (p *Program) UniformLight(name string, l ILight) {
	p.UniformColor(name + ".ambient"  , l.Ambient())
	p.UniformColor(name + ".diffuse"  , l.Diffuse())
	p.UniformColor(name + ".specular" , l.Specular())
	p.UniformVec3 (name + ".position" , l.Position())
	p.UniformVec3 (name + ".halfVector" , l.HalfVector())
	p.UniformVec3 (name + ".spotDirection" , l.SpotDirection())
	p.UniformFloat(name + ".spotExponent", l.SpotExponent())
	p.UniformFloat(name + ".spotCutoff", l.SpotCutoff())
	p.UniformFloat(name + ".spotCosCutoff", l.SpotCosCutoff())
	p.UniformFloat(name + ".constantAttenuation", l.ConstantAttenuation())
	p.UniformFloat(name + ".linearAttenuation", l.LinearAttenuation())
	p.UniformFloat(name + ".quadraticAttenuation", l.QuadraticAttenuation())
	p.UniformBool (name + ".enabled", l.Enabled())
	p.UniformBool (name + ".positional", l.Positional())
}
func (p *Program) UniformCamera(name string, c ICamera) {
	p.UniformVec3(name + ".pos" , c.Position ())
	p.UniformVec3(name + ".dir" , c.Direction())
	p.UniformVec3(name + ".up"  , c.UpVector ())
}
func (p *Program) UniformColor3(name string, c *color.Color) {
	loc := p.GetUniformLocation(name)
	loc.Uniform3f(float32(c.R), float32(c.G), float32(c.B))
}
//requires change to gl.go:
/*
func (location UniformLocation) UniformMatrix3fv(v []float32, num int, transpose bool) {
	_, p := GetGLenumType(v)
	itranspose := 0
	if transpose {itranspose = 1}
	C.glUniformMatrix3fv(C.GLint(location), C.GLsizei(1), C.GLboolean(itranspose), (*C.GLfloat)(p));
}
func (location UniformLocation) UniformMatrix4fv(v []float32, num int, transpose bool) {
	_, p := GetGLenumType(v)
	itranspose := 0
	if transpose {itranspose = 1}
	C.glUniformMatrix4fv(C.GLint(location), C.GLsizei(1), C.GLboolean(itranspose), (*C.GLfloat)(p));
}
*/

func (p *Program) UniformMat3(name string, m *mat3.Mat3) {
	loc := p.GetUniformLocation(name)
	loc.UniformMatrix3fv(m.Fv32(),1,false)
}
func (p *Program) UniformMat4(name string, m *mat4.Mat4) {
	loc := p.GetUniformLocation(name)
	loc.UniformMatrix4fv(m.Fv32(),1,false)
}
func (p *Program) UniformTexture(name string, t ITexture) {
	t.Id().Bind(gl.TEXTURE_2D)
	p.UniformInt(name, 0 )
}
func (p *Program) UniformTexture_(name string, t ITexture, activeTexture int) {
	gl.ActiveTexture(gl.GLenum(activeTexture))
	t.Id().Bind(gl.TEXTURE_2D)
	p.UniformInt(name, activeTexture-int(gl.TEXTURE0) )
}


//DEFAULT NAMES -----------------------------------
func (p *Program) SetModelMatrix(m *mat4.Mat4) {
	p.UniformMat4("uModelMatrix",m)
}
func (p *Program) SetViewMatrix(m *mat4.Mat4) {
	p.UniformMat4("uViewMatrix",m)
}
func (p *Program) SetProjectionMatrix(m *mat4.Mat4) {
	p.UniformMat4("uProjectionMatrix",m)
}
func (p *Program) SetModelViewMatrix(m *mat4.Mat4) {
	p.UniformMat4("uModelViewMatrix",m) 
}
func (p *Program) SetViewProjectionMatrix(m *mat4.Mat4) {
	p.UniformMat4("uViewProjectionMatrix",m)
}
func (p *Program) SetModelViewProjectionMatrix(m *mat4.Mat4) {
	p.UniformMat4("uModelViewProjectionMatrix",m)
}

func (p *Program) SetTextureMatrix(m *mat4.Mat4) {
	p.UniformMat4("uTextureMatrix",m)
}
func (p *Program) SetNormalMatrix(m *mat3.Mat3) {
	p.UniformMat3("uNormalMatrix",m)
}

func (p *Program) SetFrontMaterial(m IMaterial) {
	p.UniformMaterial("uFrontMaterial",m)
}
func (p *Program) SetBackMaterial(m IMaterial) {
	p.UniformMaterial("uBackMaterial",m)
}
func (p *Program) SetLight(l ILight) {
	id := l.Id()
	if id < 0 || id >= globals.MaxLights {
		panic("wrong light index/id to set in program")
	}
	p.UniformLight(fmt.Sprintf("uLightSource[%v]",l.Id()), l)
}
func (p *Program) SetCamera(c ICamera) {
	p.UniformCamera("uCamera", c)
}
func (p *Program) SetTexture(t ITexture) {
	p.UniformTexture("uTexture", t)
}
//-----------------------GLOBAL-----------------------
//BEFORE USE()
func SetCamera(c ICamera) {
	GCamera = c
}
//func SetViewMatrix(m *mat4.Mat4) {
//	GViewMatrix.Set(m)
//}
func SetProjectionMatrix(m *mat4.Mat4) {
	GProjectionMatrix.Set(m)
}
func SetFarClipplane(far Double) {
	GFarClipplane = far
}
//AFTER USE()
func SetModelMatrix(m *mat4.Mat4) {
	GUsedProgram.SetModelMatrix(m)
	GUsedProgram.SetNormalMatrix(m.Mat3().Transpose())
	m = m.Mul(GViewMatrix)
	GUsedProgram.SetModelViewMatrix(m)
	m = m.Mul(GProjectionMatrix)
	GUsedProgram.SetModelViewProjectionMatrix(m)
}
func SetTextureMatrix(m *mat4.Mat4) {
	GUsedProgram.SetTextureMatrix(m)
}
func SetFrontMaterial(m IMaterial) {
	GUsedProgram.SetFrontMaterial(m)
}
func SetBackMaterial(m IMaterial) {
	GUsedProgram.SetBackMaterial(m)
}
func SetLight(l ILight) {
	GUsedProgram.SetLight(l)
}
func SetTexture(t ITexture) {
	GUsedProgram.SetTexture(t)
}
//func SetNormalMatrix(m *mat4.Mat4) {
//	GUsedProgram.SetNormalMatrix(m)
//}

