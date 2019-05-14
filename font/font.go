package font

import (
	"fmt"
	"ogo/texture"
	. "ogo/common"
)

type Font struct {
	data *FontData
}

func NewFont(td *FontData) *Font {
	return &Font{data: td}
}

func (t *Font) Data () *FontData {return t.data}
func (t *Font) LoadTexture(text string, size int) (*texture.Texture,Double) {
	return t.data.LoadTexture(text, size)
}
func (t *Font) GetTexture(text string, size int) (*texture.Texture,Double) {
	return t.data.GetTexture(text, size)
}
func (t *Font) DeleteTexture(text string, size int) {
	t.data.DeleteTexture(text, size)
}

func (t *Font) String() string {
	return fmt.Sprintf("%v", t.data.String())
}
