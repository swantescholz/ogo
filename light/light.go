package light

import (
	"github.com/banthar/gl"
	. "ogo/glmath/color"
	. "ogo/glmath/vec3"
	"ogo/globals"
	"ogo/program"
	"ogo/glmath/material"
	. "ogo/common"
)

var (
	GFreeIds = []bool{true,true,true,true,true,true,true,true}
	GLightModelAmbient = Col(.2,.2,.2,0.0)
	GLightModelLocalViewer, GLightModelTwoSide = true, false
)

type Light struct {
	Amb, Diff, Spec *Color
	//Pos4  *Color
	//HalfVector *Color
	Pos, Dir *Vec3
	Kc, Kl, Kq Double //attentuation factors (constant, linear, quadratic)
	Se, Sc, Scc Double //spotlight exponent/cutoff/cos(cutoff)
	id int
	enabled, positional bool
}

//satisfy program.ILight interface:
func (t *Light) Id() int {return t.id}
func (t *Light) Ambient() *Color {return t.Amb}
func (t *Light) Diffuse() *Color {return t.Diff}
func (t *Light) Specular() *Color {return t.Spec}
func (t *Light) Material() *material.Material {return material.New(t.Amb,t.Diff,t.Spec,0.5)}
func (t *Light) Position() *Vec3 {return t.Pos}
func (t *Light) HalfVector() *Vec3 {return V3(1,0,0)}
func (t *Light) SpotDirection() *Vec3 {return t.Dir}
func (t *Light) SpotExponent() Double {return t.Se}
func (t *Light) SpotCutoff() Double {return t.Sc}
func (t *Light) SpotCosCutoff() Double {return t.Scc}
func (t *Light) ConstantAttenuation() Double {return t.Kc}
func (t *Light) LinearAttenuation() Double {return t.Kl}
func (t *Light) QuadraticAttenuation() Double {return t.Kq}
func (t *Light) Enabled() bool {return t.enabled}
func (t *Light) Positional() bool {return t.positional}

func New() *Light {
	l := new(Light)
	l.Amb  = White.Copy()
	l.Diff = White.Copy()
	l.Spec = White.Copy()
	l.Pos, l.Dir = V3(0,0,0), V3(0,0,-1)
	l.Kc, l.Kl, l.Kq = 1.0, 0.0, 0.0
	l.Se, l.Sc = 0.0, 180.0
	l.id = -1
	l.enabled, l.positional = false, true
	return l
}

func (t *Light) Enable() {
	for t.id = 0; !GFreeIds[t.id]; t.id++ {
		if t.id >= 7 {
			panic("too many lights")
		}
	}
	gl.Enable(getLightName(t.id))
	GFreeIds[t.id] = false
	t.enabled = true
}
func (t *Light) Disable() {
	gl.Disable(getLightName(t.id))
	GFreeIds[t.id] = true
	t.enabled = false
	t.id = -1
}

func (t *Light) Move(v *Vec3) {
	t.Pos.Addi(v)
}

//AFTER program.Use()!
func (t *Light) Shine() {
	if !t.enabled {return}
	
	t.Scc = t.Sc.Cos()
	
	if globals.UseShader {
		program.SetLight(t)
	} else {
		eid := getLightName(t.id)
		
		//light character
		gl.Lightfv(eid, gl.AMBIENT , t.Amb.Slc32());
		gl.Lightfv(eid, gl.DIFFUSE , t.Diff.Slc32());
		gl.Lightfv(eid, gl.SPECULAR, t.Spec.Slc32());
		
		//position
		var w Double = 0.0
		if t.positional {w = 1.0}
		pos4 := Col(t.Pos.X,t.Pos.Y,t.Pos.Z,w)
		gl.Lightfv(eid, gl.POSITION, pos4.Slc32());
		
		//attentuation
		gl.Lightf(eid, gl.CONSTANT_ATTENUATION , float32(t.Kc));
		gl.Lightf(eid, gl.LINEAR_ATTENUATION   , float32(t.Kl));
		gl.Lightf(eid, gl.QUADRATIC_ATTENUATION, float32(t.Kq));
		
		//spot
		gl.Lightf (eid, gl.SPOT_EXPONENT , float32(t.Se));
		gl.Lightf (eid, gl.SPOT_CUTOFF   , float32(t.Sc));
		gl.Lightfv(eid, gl.SPOT_DIRECTION, t.Dir.Slc32());
	}
}
func (t *Light) SetMaterial(m *material.Material) {
	t.Amb, t.Diff, t.Spec = m.Colors()
}
func (t *Light) SetPositional(positional bool) {
	t.positional = positional
}

func SetLightModelAmbient(amb *Color) {
	GLightModelAmbient.Set(amb)
}

func getLightName(lightID int) gl.GLenum {
	switch lightID {
		case 0: return gl.LIGHT0
		case 1: return gl.LIGHT1
		case 2: return gl.LIGHT2
		case 3: return gl.LIGHT3
		case 4: return gl.LIGHT4
		case 5: return gl.LIGHT5
		case 6: return gl.LIGHT6
		case 7: return gl.LIGHT7
		default: panic("illegal light ID")
	}
	return gl.LIGHT0
}

//activate global light model settings (amient light, local view, two sided)
func ApplyLightModelToScene() {
	if globals.UseShader {
		program.GUsedProgram.UniformColor("uLightModelAmbient"    , GLightModelAmbient)
		program.GUsedProgram.UniformBool ("uLightModelLocalViewer", GLightModelLocalViewer)
		program.GUsedProgram.UniformBool ("uLightModelTwoSide"    , GLightModelTwoSide)
	} else {
		gl.LightModelfv(gl.LIGHT_MODEL_AMBIENT, GLightModelAmbient.Slc32());
		ilocview, itwosides := 0,0
		if GLightModelLocalViewer {ilocview  = 1}
		if GLightModelTwoSide     {itwosides = 1}
		gl.LightModeli (gl.LIGHT_MODEL_LOCAL_VIEWER, ilocview);
		gl.LightModeli (gl.LIGHT_MODEL_TWO_SIDE    , itwosides);
	}
}

