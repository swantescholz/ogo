package texture

import (
	//"github.com/banthar/gl"
	"fmt"
	//"math"
	. "glmath/color"
	"image"
	_ "image/png"
	_ "image/jpeg"
	_ "image/gif"
	"os"
	"log"
)

type PixelGrid struct {
	Pix []uint8
	w, h int
}

func (t *PixelGrid) String() string {
	return fmt.Sprintf("%v, %v", t.w, t.h)
}
func (t *PixelGrid) Size () int {return t.w * t.h * 4}
func (t *PixelGrid) Sizes() (int, int) {return t.w, t.h}
func (t *PixelGrid) SetPixel(x,y int, c Color32) {
	index := (y*t.w + x)*4
	t.Pix[index+0] = c.R
	t.Pix[index+1] = c.G
	t.Pix[index+2] = c.B
	t.Pix[index+3] = c.A
}
func (t *PixelGrid) GetPixel(x,y int) (col Color32) {
	index := (y*t.w + x)*4
	col.R = t.Pix[index+0]
	col.G = t.Pix[index+1]
	col.B = t.Pix[index+2]
	col.A = t.Pix[index+3]
	return
}
func NewPixelGrid(w, h int) *PixelGrid {
	grid := &PixelGrid{make([]uint8, w*h*4), w, h}
	return grid
}

func NewPixelGridFromFile(path string) *PixelGrid {
	file, err := os.Open(path)
	if err != nil {log.Fatal(err)}
	defer file.Close()
	m, _, err := image.Decode(file)
	if err != nil {log.Fatal(err)}
	bounds := m.Bounds()
	
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y
	grid := NewPixelGrid(w,h)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := m.At(x, bounds.Max.Y - y - 1).RGBA()
			xx, yy := x - bounds.Min.X, y - bounds.Min.Y
			grid.SetPixel(xx, yy, Color32{uint8(r/256),uint8(g/256),uint8(b/256),uint8(a/256)})
		}
	}
	return grid
}




