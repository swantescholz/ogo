package texture

import (
	"github.com/banthar/gl"
	"fmt"
	. "ogo/glmath/color"
	"ogo/glmath/material"
)

type Binder interface {
	Bind()
	Unbind()
}

type Texture struct {
	data *TextureData
	//mat mat4.Mat4
}

func NewTexture(td *TextureData) *Texture {
	return &Texture{data: td}
}

func (t *Texture) Id() gl.Texture {return t.data.Id()}
func (t *Texture) Data () *TextureData {return t.data}

//AFTER program.Use()!
func (t *Texture) Bind() {
	gl.Color4d(1, 1, 1, 1)
	t.data.Bind()
}
func (t *Texture) Bind2(col *Color) {
	gl.Color4dv(col.Slc())
	t.data.Bind()
}
func (t *Texture) BindForSkybox(col *Color) {
	gl.Color4dv(col.Slc())
	t.data.Bind()
}
func (t *Texture) Bind3(mat *material.Material) {
	mat.Use()
	t.data.Bind()
}
func (t *Texture) Unbind() {
	t.data.Unbind()
}
func (t *Texture) String() string {
	return fmt.Sprintf("%v", t.data.String())
}



