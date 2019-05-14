package window

/*
MANUAL:
you must use gl.Init() for GLEW initialization -> for VBOs and shaders
if you use multithreading you must lock this thread with runtime.LockOSThread()
no other gorouting than the main goroutine may call gl functions!
if gl.GenBuffer() fails, use glGenBuffersARB() //not good!
INTEL SUCKS, check glxinfo for your GL-version
*/

import (
	"github.com/banthar/gl"
	"github.com/banthar/Go-SDL/sdl"
	. "fmt"
	"runtime"
	"time"
	"ogo/glmath/mat4"
	"ogo/util"
	"ogo/program"
	"ogo/globals"
	. "ogo/common"
)

var GWindow *Window = nil

type Window struct {
	width, height int
	title, tabname string
	projectionMatrix *mat4.Mat4
	pixelData []uint8 //used for screenshots
	elapsed Double
	lastTime time.Time
	frame int
}
func Open(w,h int, fullscreen bool, title,tabname string) *Window {
	
	if GWindow != nil {
		panic("only one Window is allowed")
	}
	win := &Window{width: w, height: h, title: title, tabname: tabname}
	GWindow = win
	
	sdl.Init(sdl.INIT_EVERYTHING)
	gl.Init() 
	
	var screen = sdl.SetVideoMode(win.width, win.height, 16, sdl.OPENGL)
	
	if screen == nil {
		sdl.Quit()
		panic("Couldn't set sdl GL video mode: " + sdl.GetError() + "\n")
	}
	sdl.WM_SetCaption(win.title, win.tabname)
	win.initGL()
	
	win.pixelData = make([]uint8, win.width*win.height*4)
	
	return win
}
func (this *Window) String() string   {return Sprintf("%v, %v", this.width, this.height)}
func (this *Window) Width () int      {return this.width}
func (this *Window) Height() int      {return this.height}
func (this *Window) ProjectionMatrix() *mat4.Mat4 {return this.projectionMatrix}
func (this *Window) Close () {
	this.pixelData = nil
	sdl.Quit()
}

func (this *Window) Elapsed() Double {return this.elapsed}
func (this *Window) FinishFrame() {
	this.flip()
	this.clearBuffers()
	this.resetModelViewMatrix()
	this.frame++
	var now = time.Now()
	if this.frame >= 2 {
		this.elapsed = Double(now.Sub(this.lastTime).Seconds())
	}
	this.lastTime = now
}
func (this *Window) ResetElapsed() {
	this.elapsed = 0.0
	this.lastTime = time.Now()
}
func (this *Window) flip  () {sdl.GL_SwapBuffers()}
func (this *Window) resetModelViewMatrix() {
	if globals.UseShader {
		program.SetModelMatrix(mat4.Identity())
	}
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}
func (this *Window) clearBuffers () {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT )
}
func (this *Window) initGL() {

	runtime.LockOSThread()
	
	if gl.Init() != 0 {
		panic("gl init error")	
	}
	
	gl.ShadeModel(gl.SMOOTH)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.LIGHTING)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.TEXTURE_2D)
	gl.TexEnvf(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.MODULATE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	var magFilter, minFilter = gl.LINEAR, gl.LINEAR
	if globals.CreateMipmaps {
		minFilter = gl.LINEAR_MIPMAP_LINEAR
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, magFilter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, minFilter)
	gl.Enable(gl.NORMALIZE)
	gl.ClearColor(.1,.1,.1,1)
	this.setGLViewport()
}
func (this *Window) setGLViewport() {
	this.reshape(this.width, this.height)
}

func (this *Window) SetProjection2D() {
	gl.MatrixMode(gl.PROJECTION)
	gl.Viewport(0, 0, int(this.Width()), int(this.Height()))
	gl.LoadIdentity()
	gl.Ortho(0, float64(this.Width()), 0, float64(this.Height()), -1.0, 1.0)
	gl.MatrixMode(gl.MODELVIEW)
}

func (this *Window) reshape(w,h int) {
	gl.Viewport(0,0,w,h)
	ratio := Double(w)/Double(h)
	var znear, zfar, fovx Double = 0.01, 65536.0, 80.0
	this.projectionMatrix = mat4.Projection(fovx.ToRad(), ratio, znear, zfar)
	if globals.UseShader {
		program.SetProjectionMatrix(this.projectionMatrix)
		program.SetFarClipplane(zfar)
	}
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.LoadMatrixd(this.projectionMatrix.Ptr())
	gl.MatrixMode(gl.MODELVIEW)
}

func (this *Window) SaveScreenshot(name string) {
	gl.ReadBuffer(gl.FRONT)
	gl.ReadPixels(0, 0, this.width, this.height, gl.RGBA, this.pixelData)
	var im = util.PixelDataToImage(this.pixelData, this.width, this.height)
	util.WritePng(im, "res/screenshots/" + name)
}

