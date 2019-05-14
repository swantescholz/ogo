package font

import (
	"log"
)


var GFontPath1, GFontPath2 = "res/fonts/", ".ttf"
func SetDefaultFontPaths(start, end string) {
	GFontPath1, GFontPath2 = start, end
}
func FontStringToPath(s string) string {return GFontPath1 + s + GFontPath2}

var GFontManager = &FontManager{make(map[string]*FontData)}
type FontManager struct {
	mfonts map[string]*FontData
}
func (t *FontManager) ClearAllData() {
	for k,v := range t.mfonts {
		v.Destroy()
		delete(t.mfonts, k)
	}
}
func (t *FontManager) Delete(path string) {
	path = FontStringToPath(path)
	v, ok := t.mfonts[path]
	if !ok {return}
	v.Destroy() //cleaning gl-memory
	delete(t.mfonts, path)
}
func (t *FontManager) Load(path string) *Font {
	path = FontStringToPath(path)
	v,ok := t.mfonts[path]
	if ok {return NewFont(v)}
	t.mfonts[path] = NewFontData(path)
	return NewFont(t.mfonts[path])
}
func (t *FontManager) Get(path string) *Font {
	path = FontStringToPath(path)
	v,ok := t.mfonts[path]
	if !ok {
		log.Printf("warning: font %v not yet loaded\n", path)
		return nil
	}
	return NewFont(v)
}

func ClearAllData()      {GFontManager.ClearAllData()}
func Delete(path string) {GFontManager.Delete(path)}
func Load  (path string) *Font {return GFontManager.Load(path)}
func Get   (path string) *Font {return GFontManager.Get(path)}


