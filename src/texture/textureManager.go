package texture

import (
	//"github.com/banthar/gl"
	//"fmt"
	//"math"
	//. "glmath/color"
	//"glmath/mat4"
	//. "glmath/util"
)

var GTexturePath1, GTexturePath2 = "res/textures/", ".png"
func SetDefaultTexturePaths(start, end string) {
	GTexturePath1, GTexturePath2 = start, end
}
func textureStringToPath(s string) string {return GTexturePath1 + s + GTexturePath2}

var GTextureManager = &TextureManager{make(map[string]*TextureData)}
type TextureManager struct {
	mtex map[string]*TextureData
}
func (t *TextureManager) ClearAllData() {
	for k,v := range t.mtex {
		v.Destroy()
		delete(t.mtex, k)
	}
}
func (t *TextureManager) Delete(path string) {
	path = textureStringToPath(path)
	v, ok := t.mtex[path]
	if !ok {return}
	v.Destroy() //cleaning gl-memory
	delete(t.mtex, path)
}
func (t *TextureManager) Load(path string) *Texture {
	path = textureStringToPath(path)
	v,ok := t.mtex[path]
	if ok {return NewTexture(v)}
	t.mtex[path] = NewTextureData(path)
	return NewTexture(t.mtex[path])
}
func (t *TextureManager) Get(path string) *Texture {
	path = textureStringToPath(path)
	v,ok := t.mtex[path]
	if !ok {return nil}
	return NewTexture(v)
}

func ClearAllData()      {GTextureManager.ClearAllData()}
func Delete(path string) {GTextureManager.Delete(path)}
func Load  (path string) *Texture {return GTextureManager.Load(path)}
func Get   (path string) *Texture {return GTextureManager.Get(path)}


