package texture

import (
	"github.com/banthar/gl"
	"fmt"
	//"math"
	. "glmath/color"
	//"glmath/mat4"
	//. "glmath/util"
)

type Texture struct {
	data *TextureData
	//mat mat4.Mat4
}

func NewTexture(td *TextureData) *Texture {
	return &Texture{data: td}
}

func (t *Texture) Data () *TextureData {return t.data}
func (t *Texture) Bind() {
	gl.Color4d(1, 1, 1, 1)
	t.data.Bind()
}
func (t *Texture) Bind2(col *Color) {
	gl.Color4d(col.R, col.G, col.B, col.A)
	t.data.Bind()
}
func (t *Texture) Unbind() {
	t.data.Unbind()
}
func (t *Texture) String() string {
	return fmt.Sprintf("%v", t.data.String())
}



