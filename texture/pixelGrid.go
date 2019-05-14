package texture

import (
	"fmt"
	"image"
	"os"
	"log"
	. "ogo/glmath/color24"
	. "ogo/glmath/color32"
	"ogo/util"
)

type PixelGrid struct {
	Pix []uint8
	w, h int
}

var index int

func (t *PixelGrid) String() string {
	return fmt.Sprintf("%v, %v", t.w, t.h)
}
func (t *PixelGrid) Size () int {return t.w * t.h * 4}
func (t *PixelGrid) Sizes() (int, int) {return t.w, t.h}
func (t *PixelGrid) SetPixel(x,y int, col *Color32) {
	index = t.coordsToIndex(x,y)
	t.Pix[index+0] = col.R
	t.Pix[index+1] = col.G
	t.Pix[index+2] = col.B
	t.Pix[index+3] = col.A
}
func (t *PixelGrid) SetPixel24(x,y int, col *Color24) {
	index = t.coordsToIndex(x,y)
	t.Pix[index+0] = col.R
	t.Pix[index+1] = col.G
	t.Pix[index+2] = col.B
}
func (t *PixelGrid) GetPixel(x,y int) *Color32 {
	index = t.coordsToIndex(x,y)
	var col = new(Color32)
	col.R = t.Pix[index+0]
	col.G = t.Pix[index+1]
	col.B = t.Pix[index+2]
	col.A = t.Pix[index+3]
	return col
}
func (t *PixelGrid) GetPixel24(x,y int) *Color24 {
	index = t.coordsToIndex(x,y)
	var col = new(Color24)
	col.R = t.Pix[index+0]
	col.G = t.Pix[index+1]
	col.B = t.Pix[index+2]
	return col
}
func (t *PixelGrid) SetAlpha(a uint8) {
	for y := 0; y < t.h; y++ {
		for x := 0; x < t.w; x++ { 
			t.Pix[t.coordsToIndex(x,y)+3] = a
		}	
	}
}

func (t *PixelGrid) coordsToIndex(x,y int) int {return (y*t.w + x)*4}

func (t *PixelGrid) Image() image.Image {
	return util.PixelDataToImage(t.Pix,t.w,t.h)
}
func (t *PixelGrid) WriteToFile(path string) {
	util.WritePng(t.Image(), path)
}

func (t *PixelGrid) CreateMipmap(w,h int) *PixelGrid {
	var facw, fach = float64(t.w)/float64(w), float64(t.h)/float64(h)
	var pix = NewPixelGrid(w,h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var col = t.GetPixel(int(float64(x)*facw), int(float64(y)*fach))
			pix.SetPixel(x,y, col)
		}
	}
	return pix
}

//---------------

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
			grid.SetPixel(xx, yy, NewColor32(uint8(r/256),uint8(g/256),uint8(b/256),uint8(a/256)))
		}
	}
	return grid
}




