package texture

import (
	"github.com/banthar/gl"
	"fmt"
	//"math"
	. "glmath/color"
	//"glmath/mat4"
	//. "glmath/util"
)


type TextureData struct {
	id gl.Texture
	path string
	//stdMat mat4.Mat4
	Pix *PixelGrid
	Mips []*PixelGrid
}

func NewTextureData(path string) *TextureData {
	id := gl.GenTexture()
	tb := &TextureData{id, path, nil, make([]*PixelGrid,0)}
	tb.Pix = NewPixelGridFromFile(tb.path)
	tb.CreateMipmaps()
	tb.LoadIntoGL()
	return tb
}

//clears the opengl texture memory
func (t *TextureData) Destroy() {
	t.id.Delete()
}

func (t *TextureData) Id    () gl.Texture {return t.id}
func (t *TextureData) Name  () string {return t.path}
func (t *TextureData) Sizes () (int, int) {return t.Pix.Sizes()}
func (t *TextureData) Width () int {w,_ := t.Pix.Sizes(); return w}
func (t *TextureData) Height() int {_,h := t.Pix.Sizes(); return h}

func (t *TextureData) CreateMipmaps() {
	w, h := t.Width()/2, t.Height()/2
	var pix *PixelGrid = nil
	lastMip := t.Pix
	for w > 0 && h > 0 {
		pix = NewPixelGrid(w,h)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				var col Color32 = lastMip.GetPixel(x*2, y*2)
				pix.SetPixel(x,y, col)
				//pix.SetPixel(x,y, Color32{255,uint8(x-y),0,255})
			}
		}
		t.Mips = append(t.Mips, pix)
		w, h = w/2, h/2
		lastMip = pix
	}
}
func (t *TextureData) LoadIntoGL() {
	t.id.Bind(gl.TEXTURE_2D)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, t.Width(), t.Height(), 0, gl.RGBA, gl.UNSIGNED_BYTE, t.Pix.Pix)
	for i,v := range t.Mips {
		gl.TexImage2D(gl.TEXTURE_2D, i+1, gl.RGBA, v.w, v.h, 0, gl.RGBA, gl.UNSIGNED_BYTE, v.Pix)
	}
}

func (t *TextureData) Bind() {
	t.id.Bind(gl.TEXTURE_2D)
}
func (t *TextureData) Unbind() {
	t.id.Unbind(gl.TEXTURE_2D)
}
func (t *TextureData) String() string {
	return fmt.Sprintf("%v, %v", t.id, t.path)
}



