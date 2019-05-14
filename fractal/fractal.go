package fractal

import (
	. "github.com/salviati/gmp"
	. "fmt"
	"io/ioutil"
	"strings"
	. "ogo/glmath/color24"
	. "ogo/glmath/vec2"
	. "ogo/common"
	"ogo/texture"
	"ogo/util"
	. "ogo/glmath/gradient24"
	"ogo/ass"
	//*
	cx "ogo/glmath/cx/oldcx"/*///
	"glmath/cx"//*/
)

const (
	configfile = "res/fractalConfig.txt"
)

type Fractal struct {
	width, height int
	gradient *Gradient
	Iterations, Precision int
	Center, Z0 *cx.Cx
	Zoom *cx.Cx
	tex *texture.Texture
	grid *texture.PixelGrid
}

func New(gradient *Gradient) *Fractal {
	var f = new(Fractal)
	f.width = 256
	f.height = 256
	f.gradient = gradient
	f.Iterations = 12
	f.Center = cx.New(0,0)
	f.Z0 = cx.New(0,0)
	f.Zoom = cx.New(3.5,1.0)
	f.Precision = 64
	return f
}

func (f *Fractal) Quit() {
	b, e := ioutil.ReadFile(configfile)
	ass.Error(e)
	lines := strings.Split(string(b), "\n")
	var index = 0
	var next = func(s string) {
		lines[index] = s
		index += 1
	}
	next(Sprintf("%v %v",f.width, f.height))
	next(Sprintf("%v %v",f.Iterations, f.Precision))
	next(Sprint(f.Center.X))
	next(Sprint(f.Center.Y))
	next(Sprint(f.Zoom.X))
	ioutil.WriteFile(configfile, []byte(strings.Join(lines,"\n")), 0)
	
	f.Center.Clear()
	f.Z0.Clear()
	f.Zoom.Clear()
}

func (f *Fractal) String() string {
	return Sprintf("FRACTAL: %v (%v %v) %v, %v, %v, %v", cx.FloatType, f.width,f.height,f.Center, f.Zoom.X, f.Iterations, f.Precision)
}

func (f *Fractal) Width() int {return f.width}
func (f *Fractal) Height() int {return f.height}
func (f *Fractal) Ratio() Double {return Double(f.height) / Double(f.width)}
func (f *Fractal) Texture() *texture.Texture {return f.tex}
func (f *Fractal) Gradient() *Gradient {return f.gradient}

func (f *Fractal) reloadPrecision() {
	SetDefaultPrec(f.Precision)
	cx.SetPrec(f.Precision)
	f.Zoom.SetPrec(f.Precision)
	f.Z0.SetPrec(f.Precision)
	f.Center.SetPrec(f.Precision)
}

func (f *Fractal) CreateTexture() *texture.Texture {
	f.reloadPrecision()
	var ratio = cx.NewFloat(f.Ratio())
	f.Zoom.SetY(f.Zoom.X)
	f.Zoom.MulY(ratio)
	var halfToUpperRight = f.Zoom.MulsDouble(.5)
	var lowerLeft = f.Center.Sub(halfToUpperRight)
	var wh = cx.New(Double(f.width), Double(f.height))
	var fxy = f.Zoom.DivComponents( wh )
	var coord = lowerLeft.Copy()
	defer func() {
		ratio.Clear()
		halfToUpperRight.Clear(); lowerLeft.Clear()
		wh.Clear(); fxy.Clear()
		coord.Clear()
	}()
	var sum Double = 0.0
	var col *Color24 = NewColor24(0,0,0)
	var intensity Double
	Printf("Idle time: "); util.Ts()
	cx.SetIterations(f.Iterations)
	for y := 0; y < f.height; y++ {
		coord.SetX(lowerLeft.X)
		coord.AddY(fxy.Y)
		for x := 0; x < f.width; x++ {
			coord.AddX(fxy.X)
			intensity = Double(f.Z0.DoMandel(coord)) / Double(f.Iterations)
			//intensity = intensity.Sqrt()
			f.gradient.Interpolate(col, intensity)
			f.grid.SetPixel24(x, y, col)
		}
		sum += intensity
	}
	Printf(">>>>>>>>>>>>>>>> Math time: "); util.Ts()
	f.renewTexture()
	Printf("GL time: "); util.Ts()
	Println("Intensity sum:", sum,col)
	return f.tex
}

//move vector in texcoords
func (f *Fractal) MoveCenter(v *Vec2) {
	var c = cx.NewFromVec2(v)
	var pointFive = cx.New(.5,.5)
	defer c.Clear()
	defer pointFive.Clear()
	c.Subi(pointFive)
	c.MuliComponents(f.Zoom)
	f.Center.Addi(c)
}

func (f *Fractal) ZoomIn(factor Double) {
	var fzoom = cx.NewFloat(factor)
	defer fzoom.Clear()
	f.Zoom.Mulsi(fzoom)
}

func rangewh(w,h int, f func(x,y int)) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			f(x,y)
		}
	}
}

func (f *Fractal) renewTexture() {
	f.tex = texture.NewTexture(texture.NewTextureDataFromPixelGrid(f.grid))
}
func (f *Fractal) RenewResolution(w,h int) {
	f.width, f.height = w,h
	f.height = f.width
	var last = f.grid
	f.grid = texture.NewPixelGrid(f.width, f.height)
	if last != nil {
		w,h = last.Sizes()
		var xmax,ymax = f.grid.Sizes()
		rangewh(w,h, func(x,y int) {
			if x < xmax && y < ymax {
				var col = last.GetPixel24(x,y)
				f.grid.SetPixel24(x,y,col)
			}			
		})
	}
	f.grid.SetAlpha(255)
	f.renewTexture()
}

func (f *Fractal) LoadConfig() {
	b, e := ioutil.ReadFile(configfile)
	ass.Error(e)
	lines := strings.Split(string(b), "\n")
	var index = 0
	var next = func() string {
		defer func(){index += 1}()
		return lines[index]
	}
	Sscanf(next(), "%v %v", &f.width, &f.height)
	f.height = f.width
	Sscanf(next(), "%v %v", &f.Iterations, &f.Precision)
	f.reloadPrecision()
	var x,y = next(),next()
	Println(x,y)
	f.Center.SetString(x,y)
	f.Zoom.SetString(next(), "1.0")
	f.RenewResolution(f.width, f.height)
}

