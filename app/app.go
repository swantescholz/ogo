package app

import (
	"github.com/banthar/Go-SDL/sdl"
	"github.com/banthar/gl"
	"github.com/salviati/gmp"
	. "fmt"
	"ogo/util"
	"time"
	"ogo/window"
	"io/ioutil"
	. "ogo/glmath/vec2"
	. "ogo/glmath/color"
	"ogo/glmath/color24"
	. "ogo/glmath/gradient24"
	"strings"
	. "ogo/common"
	"ogo/input/listener"
	. "ogo/input/constants"
	"ogo/texture"
	"ogo/fractal"
)

var T0 uint32 = 0
var Frames uint32 = 0

type App struct {
	win *window.Window
	il *listener.Listener
	elapsed Double
	quit bool
	frac *fractal.Fractal
	texmandel *texture.Texture
	cursor *Cursor
	moveSpeed, zoomStep Double
	iterationStep, precStep int
}

type Cursor struct {
	Pos *Vec2 //in tex coords
	Zoom Double
}

func NewCursor(pos *Vec2, zoom Double) *Cursor {
	return &Cursor{pos, zoom}
}

func (this *Cursor) String() string {return Sprintf("%v, %v", this.Pos, this.Zoom)}
func (this *Cursor) Move(v *Vec2) {this.Pos.Addi(v)}
func (this *Cursor) ZoomIn(zoomStep Double) {this.Zoom = zoomStep}
func (this *Cursor) ZoomOut(zoomStep Double) {this.ZoomIn(zoomStep.Inv())}
func (this *Cursor) Draw(w,h Double, col *color24.Color24) {
	var p = this.Pos.Mul(V2(w,h))
	var a = p.Add(V2(-10,-10))
	var b = p.Add(V2(10,-10))
	gl.Disable(gl.TEXTURE_2D)
	col.Gl()
	gl.Begin(gl.TRIANGLES)
	p.Gl(); a.Gl(); b.Gl()
	gl.End()
	gl.Enable(gl.TEXTURE_2D)
}

func New() *App {
	var this = &App{il: listener.New()}
	return this
}

func (this *App) readConfig(configfile string) (width,height int, fullscreen bool, title,tabname string) {
	b, e := ioutil.ReadFile(configfile)
	if e != nil {panic(e)}
	lines := strings.Split(string(b), "\n")
	var strfs string
	
	Sscanf(lines[0], "%v", &width)
	Sscanf(lines[1], "%v", &height)
	Sscanf(lines[2], "%v", &strfs)
	if strfs == "0" || strfs == "true" {fullscreen = true}
	Sscanf(lines[3], "%v", &title)
	Sscanf(lines[4], "%v", &tabname)
	Sscanf(lines[5], "%v", &this.moveSpeed)
	Sscanf(lines[6], "%v", &this.zoomStep)
	Sscanf(lines[7], "%v", &this.iterationStep)
	Sscanf(lines[8], "%v", &this.precStep)
	return
}

func (this *App) Run() {
	Println("START")
	time1 := time.Now()
	
	this.init_()
	this.mainLoop()
	this.quit_()
	
	dur := time.Since(time1)
	Println("END")
	Println("dt:", dur)
}

func (this *App) mainLoop() {
	for !this.quit {
		this.elapsed = this.win.Elapsed()
		this.processInput()	
		this.draw()
	}	
}

func (this *App) processFractalInput(l,r,u,d,f,b,ita,itb uint32) {
	var ms Double = this.moveSpeed * this.elapsed
	var redraw, move = false,false
	if this.il.KeyDown(Left) {this.cursor.Move(V2(-ms,0))}
	if this.il.KeyDown(Right) {this.cursor.Move(V2(ms,0))}
	if this.il.KeyDown(Up) {this.cursor.Move(V2(0,ms))}
	if this.il.KeyDown(Down) {this.cursor.Move(V2(0,-ms))}
	if this.il.KeyPressed(Q) {
		var w = this.frac.Width() / 2
		this.frac.RenewResolution(w,w)
		this.texmandel = this.frac.Texture()
		//redraw = true
	}
	if this.il.KeyPressed(A) {
		var w = this.frac.Width() * 2
		this.frac.RenewResolution(w,w)
		this.texmandel = this.frac.Texture()
		//redraw = true
	}
	if this.il.KeyPressed(R) {redraw = true}
	if this.il.KeyPressed(W) {this.frac.Precision -= this.precStep;redraw = true}
	if this.il.KeyPressed(E) {this.frac.Precision += this.precStep;redraw = true}
	if this.il.KeyPressed(C) {this.cursor.ZoomIn(this.zoomStep); move = true}
	if this.il.KeyPressed(X) {this.cursor.ZoomOut(this.zoomStep); move = true}
	if this.il.KeyPressed(S) {this.frac.Iterations -= this.iterationStep; redraw = true}
	if this.il.KeyPressed(D) {this.frac.Iterations += this.iterationStep; redraw = true}
	if redraw || move {
		if move {
			this.frac.MoveCenter(this.cursor.Pos)
			this.cursor.Pos.Set(V2(.5,.5))
			this.frac.ZoomIn(this.cursor.Zoom)
		}
		this.reloadFractal()
	}
}

func (this *App) processInput() {
	this.il.Update()

	this.processFractalInput(Left,Right,Up,Down,C,X,S,D)
	
	if this.il.KeyDown(Escape) || this.il.MbQuit() {this.quit = true}
	
	
	if this.il.KeyPressed(F1) {
		this.win.SaveScreenshot(util.NowAsString())
	}
}

func (this *App) reloadFractal() {
	Println(this.frac)
	this.texmandel = this.frac.CreateTexture()
	this.win.ResetElapsed()
}

func drawTexRect(tex *texture.Texture, a,b *Vec2) {
	tex.Bind()
	gl.Begin(gl.QUADS)
		gl.TexCoord2i(0,0); a.Gl()
		gl.TexCoord2i(1,0); V2(b.X, a.Y).Gl()
		gl.TexCoord2i(1,1); b.Gl()
		gl.TexCoord2i(0,1); V2(a.X, b.Y).Gl()
	gl.End()
}

func (this *App) draw() {
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.LIGHTING)
	gl.Disable(gl.DEPTH_TEST)
	
	White.Gl()
	var w,h = Double(this.win.Width()), Double(this.win.Height())
	drawTexRect(this.texmandel, V2(0,0), V2(w,h))
	var col = color24.NewColor24(0,0,0)
	this.frac.Gradient().Interpolate(col,this.cursor.Pos.X)
	this.cursor.Draw(w,h,col)
	
	this.win.FinishFrame()
}

func (this *App) init_() {
	
	var w,h,fullscreen,title,tabname = this.readConfig("res/config.txt")
	
	this.win = window.Open(w,h,fullscreen,title,tabname)
	this.win.SetProjection2D()
	
	var gradient = NewGradient([]*Color{Black, Blue, Red, Gold})
	this.frac = fractal.New(gradient)
	this.frac.LoadConfig()
	this.cursor = NewCursor(V2(.5,.5), 1.0)
	this.reloadFractal()
}



func (this *App) quit_() {
	this.frac.Quit()
	gmp.ClearGarbage()
	this.win.Close()
}

func (this *App) fps() {
	Frames++
	t := sdl.GetTicks()
	if t-T0 >= 5000 {
		seconds := (t - T0) / 1000.0
		fps := Frames / seconds
		print(Frames, " frames in ", seconds, " seconds = ", fps, " FPS\n")
		T0 = t
		Frames = 0
	}
}

