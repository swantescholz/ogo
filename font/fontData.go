package font

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
	. "fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"strings"
	. "ogo/glmath/color32"
	"ogo/texture"
	. "ogo/common"
)

const (
	widthPerChar = 0.7
	cspacing = 1.2
	cdpi = 72
	yfactor = cspacing
)

type FontTexture struct {
	rgba *image.RGBA
	context *freetype.Context
	textureData *texture.TextureData
	text, name string
	size int
	widthFactor Double
}

func NewFontTexture(font *truetype.Font, text string, size int) *FontTexture {
	this := new(FontTexture)
	this.text, this.size, this.name = text, size, textPlusSizeToName(text, size)
	
	lines := strings.Split(text, "\n")
	nlines := len(lines)
	maxLen := 0
	for _,s := range lines {
		if len(s) > maxLen {maxLen = len(s)}
	}
	// Initialize the context.
	fsize := Double(size)
	fg := image.White
	bg := image.Transparent
	//estimate width and height:
	width := int(widthPerChar*fsize*Double(maxLen))
	height := int(Double(nlines)*fsize*yfactor)
	Println(height, nlines)
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(cdpi)
	c.SetFont(font)
	c.SetFontSize(fsize.F())
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	
	// Draw the text.
	ptf32 := c.PointToFix32((fsize * cspacing).F())
	pt := freetype.Pt(0, size)
	//pt := freetype.Pt(0, 1)
	for _, s := range lines {
		_, err := c.DrawString(s, pt)
		//Println(s)
		if err != nil {
			log.Println(err)
			return nil
		}
		pt.Y += ptf32
	}
	
	//find actual bounds:
	xmin, xmax, ymin, ymax := 111111,0,111111,0
	bounds := rgba.Bounds()
	//Println(bounds.Min.X, bounds.Max.X, bounds.Min.Y, bounds.Max.Y)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := rgba.At(x, bounds.Max.Y - y - 1).RGBA()
			if a != 0 {
				if x < xmin {xmin = x}
				if x > xmax {xmax = x}
				if y < ymin {ymin = y}
				if y > ymax {ymax = y}
			}
		}
	}
	//Println(bounds.Min.X, bounds.Max.X, bounds.Min.Y, bounds.Max.Y)
	//fill a pixel grid for the textureData
	w,h := xmax - xmin + 1, ymax - ymin + 1
	grid := texture.NewPixelGrid(w,h)
	for y := ymin; y <= ymax; y++ {
		for x := xmin; x <= xmax; x++ {
			r, g, b, a := rgba.At(x, bounds.Max.Y - y - 1).RGBA()
			xx, yy := x - xmin, y - ymin
			grid.SetPixel(xx, yy, NewColor32(uint8(r/256),uint8(g/256),uint8(b/256),uint8(a/256)))
		}
	}
	
	grid.WriteToFile("fonttest")
	this.rgba, this.context = rgba, c
	td := texture.NewTextureData(this.name)
	td.Pix = grid
	td.CreateMipmaps()
	td.LoadIntoGL()
	this.textureData = td
	w, h = this.textureData.Sizes()
	this.widthFactor = Double(w)/Double(h)
	return this
}
func (this *FontTexture) Destroy() {
	this.textureData.Destroy()
}
func (this *FontTexture) CreateGlTexture() *texture.Texture {
	return texture.NewTexture(this.textureData)
}

type FontData struct {
	font *truetype.Font
	name string
	mtex map[string]*FontTexture
}

func NewFontData(fontpath string) *FontData {
	// Read the font data.
	fontBytes, err := ioutil.ReadFile(fontpath)
	if err != nil {
		log.Fatal(err)
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatal(err)
	}
	return &FontData{font: font, name: fontpath, mtex: make(map[string]*FontTexture)}
}

func (this *FontData) String() string {
	return Sprintf("%v", this.name)
}

func (this *FontData) Destroy() {
	this.ClearTextures()
}
func (this *FontData) ClearTextures() {
	for k,v := range this.mtex {
		v.Destroy()
		delete(this.mtex, k)
	}
}

func (this *FontData) Name  () string {return this.name}

func textPlusSizeToName(text string, size int) string {
	return Sprintf("%v ::%v", size, text)
}
func nameToTextPlusSize(name string) (text string, size int) {
	Sscanf(name[:3], "%v", &size)
	text = name[6:]
	return
}

func (this *FontData) LoadTexture(text string, size int) (gltexture *texture.Texture, widthFactor Double) {
	name := textPlusSizeToName(text, size)
	if !this.textureNameExists(name) {
		this.mtex[name] = NewFontTexture(this.font, text, size)
	}
	return this.GetTexture(text, size)
}
func (this *FontData) GetTexture(text string, size int) (gltexture *texture.Texture, widthFactor Double) {
	name := textPlusSizeToName(text, size)
	if !this.textureNameExists(name) {
		panic("font texture to get does not exist")
	}
	var fonttex = this.mtex[name]
	gltexture = fonttex.CreateGlTexture()
	widthFactor = fonttex.widthFactor
	return
}
func (this *FontData) TextureExists(text string, size int) bool {
	return this.textureNameExists(textPlusSizeToName(text, size))
}
func (this *FontData) textureNameExists(name string) bool {
	_,ok := this.mtex[name]
	return ok
}
func (this *FontData) DeleteTexture(text string, size int) {
	name := textPlusSizeToName(text, size)
	_,ok := this.mtex[name]
	if !ok {
		panic("font texture to be deleted does not exist")
	}
	this.mtex[name].Destroy()
	delete(this.mtex, name)
}




