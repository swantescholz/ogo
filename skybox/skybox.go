package skybox

import (
	"github.com/banthar/gl"
	. "fmt"
	"ogo/texture"
	. "ogo/glmath/vec3"
	"ogo/glmath/color"
	"ogo/globals"
	"ogo/program"
	. "ogo/common"
)

var (
	GPathStart, GPathEnd = "skybox/", ""
	GFrontName = "/ft"
	GBackName  = "/bk"
	GLeftName  = "/lt"
	GRightName = "/rt"
	GUpName    = "/up"
	GDownName  = "/dn"
	GDefaultDistance Double = 5000.0
)

type Skybox struct {
	name string
	ft, bk, lt, rt, up, dn *texture.Texture
}

//has to be called after opengl/glsl is ready!
func New(name string) *Skybox {
	s := new(Skybox)
	s.Load(name)
	return s
}

func (t *Skybox) String() string {
	return Sprintf("%v", t.name)
}

func (t *Skybox) Render(center *Vec3) {
	t.RenderInDistance(center, GDefaultDistance)
}
func (t *Skybox) RenderInDistance(center *Vec3, d Double) {
	t.RenderScaled(center, V3(1,1,1).Muls(d))
}
func (t *Skybox) RenderScaled(center, scale *Vec3) {
	v := []*Vec3{
		V3(-1,-1, 1),
		V3( 1,-1, 1),
		V3( 1, 1, 1),
		V3(-1, 1, 1),
		V3(-1,-1,-1),
		V3( 1,-1,-1),
		V3( 1, 1,-1),
		V3(-1, 1,-1)}
	for i := 0; i < 8; i++ {
		v[i].Muli(scale).Addi(center)
	}
	
	if globals.UseShader {
		program.Unuse()
	}
	
	//*//save attributes and change them
	gl.PushAttrib(gl.ENABLE_BIT | gl.TEXTURE_BIT)
	defer gl.PopAttrib() //reset to old attributes
	gl.Enable (gl.TEXTURE_2D)
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.LIGHTING)
	gl.Disable(gl.BLEND)//*/
	oneSide := func(tex *texture.Texture, a,b,c,d int, n *Vec3) {
		tex.BindForSkybox(color.White)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.Begin(gl.TRIANGLES)
			gl.Normal3dv(n.Slc())
			gl.TexCoord2i(0,0); gl.Vertex3dv(v[a].Slc())
			gl.TexCoord2i(1,0); gl.Vertex3dv(v[b].Slc())
			gl.TexCoord2i(1,1); gl.Vertex3dv(v[c].Slc())
			gl.TexCoord2i(0,0); gl.Vertex3dv(v[a].Slc())
			gl.TexCoord2i(1,1); gl.Vertex3dv(v[c].Slc())
			gl.TexCoord2i(0,1); gl.Vertex3dv(v[d].Slc())
		gl.End()
	}
	oneSide(t.up, 2,3,7,6, V3( 0,-1, 0))
	oneSide(t.dn, 5,4,0,1, V3( 0, 1, 0))
	oneSide(t.lt, 5,1,2,6, V3( 1, 0, 0))
	oneSide(t.rt, 0,4,7,3, V3(-1, 0, 0))
	oneSide(t.ft, 1,0,3,2, V3( 0, 0,-1))
	oneSide(t.bk, 4,5,6,7, V3( 0, 0, 1))
}

func (t *Skybox) Load(name string) {
	if len(t.name) > 0 {panic("skybox already loaded: " + t.name + " vs " + name)}
	t.name = name
	start := GPathStart + t.name
	end := GPathEnd
	t.ft = texture.Load(start + GFrontName + end)
	t.bk = texture.Load(start + GBackName  + end)
	t.lt = texture.Load(start + GLeftName  + end)
	t.rt = texture.Load(start + GRightName + end)
	t.up = texture.Load(start + GUpName    + end)
	t.dn = texture.Load(start + GDownName  + end)
}
func (t *Skybox) Destroy() {
	start := GPathStart + t.name
	end := GPathEnd
	texture.Delete(start + GFrontName + end)
	texture.Delete(start + GBackName  + end)
	texture.Delete(start + GLeftName  + end)
	texture.Delete(start + GRightName + end)
	texture.Delete(start + GUpName    + end)
	texture.Delete(start + GDownName  + end)
	t.ft, t.bk, t.lt, t.rt, t.up, t.dn = nil,nil,nil,nil,nil,nil
}




