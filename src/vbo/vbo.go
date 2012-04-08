package vbo

import (
	"github.com/banthar/gl"
	"fmt"
	. "glmath/vec3"
	. "glmath/vec2"
	//. "glmath/util"
	//. "glmath/color"
	//"texture"
)//*/

type Vbo struct {
	P []*Vec3
	N []*Vec3
	T []*Vec2
	I []int32
	p,n,t []float32
	idp, idn, idt, idi gl.Buffer
	created, enabled bool
	specp, specn, spect, speci bool
	usage gl.GLenum
}

func New() *Vbo {
	v := new(Vbo)
	v.InitData()
	return v
}
func (v *Vbo) String() string {
	return fmt.Sprintf("Vbo( %v, %v, %v, %v; %v, %v; %v, %v, %v, %v; %v )",
		v.idp, v.idn, v.idt, v.idi, v.created, v.enabled, v.specp, v.specn,
		v.spect, v.speci, v.usage)
}
func (v *Vbo) Destroy() {
	if v.idp != 0 {v.idp.Delete()}
	if v.idn != 0 {v.idn.Delete()}
	if v.idt != 0 {v.idt.Delete()}
	if v.idi != 0 {v.idi.Delete()}
	v.ClearData()
}
func (v *Vbo) InitData() {
	v.P,v.N,v.T,v.I = nil,nil,nil,nil
	v.idp,v.idn,v.idt,v.idi = 0,0,0,0
	v.created, v.enabled = false,false
	v.specp,v.specn,v.spect,v.speci = false,false,false,false
}
func (v *Vbo) ClearData() {
	v.InitData()
}
func (v *Vbo) Create(usage gl.GLenum) {
	v.usage = usage
	if v.created {return}
	if len(v.P) == 0 {panic("No position data specified")}
	sizep := len(v.P) * 3 * 8/2
	sizen := len(v.N) * 3 * 8/2
	sizet := len(v.T) * 2 * 8/2
	sizei := len(v.I) * 1 * 4
	//refill in simpler slices:
	
	v.p = make([]float32, 0, len(v.P)*3)
	v.n = make([]float32, 0, len(v.N)*3)
	v.t = make([]float32, 0, len(v.T)*2)
	for _,x := range v.P {v.p = append(v.p, float32(x.X), float32(x.Y), float32(x.Z))}
	for _,x := range v.N {v.n = append(v.n, float32(x.X), float32(x.Y), float32(x.Z))}
	for _,x := range v.T {v.t = append(v.t, float32(x.X), float32(x.Y))}
	if v.idp == 0 {v.idp = gl.GenBuffer()}
	if v.idn == 0 {v.idn = gl.GenBuffer()}
	if v.idt == 0 {v.idt = gl.GenBuffer()}
	if v.idi == 0 {v.idi = gl.GenBuffer()}
	if v.idp * v.idn * v.idt * v.idi == 0 {panic("Some buffer could not be generated")}
	nilSlc := []int32{0}
	if len(v.I) > 0 {
		v.idi.Bind(gl.ELEMENT_ARRAY_BUFFER)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, sizei, nilSlc, v.usage)
		gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, 0, sizei, nilSlc)
		v.speci = true
	}
	UnbindIbo()
	if len(v.P) > 0 {
		v.idp.Bind(gl.ARRAY_BUFFER)
		gl.BufferData(gl.ARRAY_BUFFER, sizep, nilSlc, v.usage)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, sizep, (&v.p[0]))
		v.specp = true
	}
	if len(v.N) > 0 {
		v.idn.Bind(gl.ARRAY_BUFFER)
		gl.BufferData(gl.ARRAY_BUFFER, sizen, nilSlc, v.usage)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, sizen, (&v.n[0]))
		v.specn = true
	}
	if len(v.T) > 0 {
		v.idt.Bind(gl.ARRAY_BUFFER)
		gl.BufferData(gl.ARRAY_BUFFER, sizet, nilSlc, v.usage)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, sizet, (&v.t[0]))
		v.spect = true
	}
	UnbindVbo()
	v.created = true
}
func (v *Vbo) Enable() {
	if v.speci {
		v.idi.Bind(gl.ELEMENT_ARRAY_BUFFER)
		//gl.IndexPointer(gl.INT, 0, 0) //TODO
		panic("index pointer unimplemented")
		gl.EnableClientState(gl.INDEX_ARRAY)
	}
	if v.specn {
		v.idn.Bind(gl.ARRAY_BUFFER)
		gl.NormalPointerVBO(gl.FLOAT, 0, 0)
		gl.EnableClientState(gl.NORMAL_ARRAY)
	}
	if v.spect {
		v.idt.Bind(gl.ARRAY_BUFFER)
		gl.TexCoordPointerVBO(2, gl.FLOAT, 0, 0)
		gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)
	}
	if v.specp {
		v.idp.Bind(gl.ARRAY_BUFFER)
		gl.VertexPointerVBO(3, gl.FLOAT, 0, 0)
		gl.EnableClientState(gl.VERTEX_ARRAY)
	}
	v.enabled = true
}
func (v *Vbo) Disable() {
	UnbindBoth()
	if v.speci {gl.DisableClientState(gl.INDEX_ARRAY)}
	if v.specn {gl.DisableClientState(gl.NORMAL_ARRAY)}
	if v.spect {gl.DisableClientState(gl.TEXTURE_COORD_ARRAY)}
	if v.specp {gl.DisableClientState(gl.VERTEX_ARRAY)}
	v.enabled = false
}


func (v *Vbo) Draw(mode gl.GLenum) {
	/*
	gl.Begin(mode)
	for i := range v.P {
		if v.specn {gl.Normal3d(v.N[i].X, v.N[i].Y, v.N[i].Z)}
		if v.spect {gl.TexCoord2d(v.T[i].X, v.T[i].Y)}
		gl.Vertex3d(v.P[i].X, v.P[i].Y, v.P[i].Z)
	}
	gl.End()
	return//*/
	
	if v.speci {
		v.DrawElements(mode, -1)
	} else {
		v.DrawArrays(mode, 0, -1)
	}
}
func (v *Vbo) DrawElements(mode gl.GLenum, count int) {
	if count < 0 {count = int(len(v.I))}
	gl.DrawElementsVBO(mode, gl.UNSIGNED_INT, count)
}
func (v *Vbo) DrawArrays(mode gl.GLenum, first, count int) {
	if count < 0 {count = int(len(v.P))}
	gl.DrawArrays(mode,first, count)
}

func UnbindVbo()  {gl.BufferUnbind(gl.ARRAY_BUFFER)}
func UnbindIbo()  {gl.BufferUnbind(gl.ELEMENT_ARRAY_BUFFER)}
func UnbindBoth() {UnbindVbo(); UnbindIbo()}




