package window

import (
	"github.com/banthar/gl"
	"github.com/jteeuwen/glfw"
	"fmt"
	"math"
	"math/rand"
	"time"
	. "glmath/util"
	"glmath/mat4"
	"input/listener"
	"camera"
)

var vFMT = fmt.Sprintf("keep 'fmt' import during debugging"); // TODO: remove 
var vRAND = rand.Float64()
var vTIME = time.NewTicker(1e9 / 2 )
var vMATH = math.Sqrt(10) //*/

var GWindow *Window = nil

type Window struct {
	width, height int
	title, tabname string
	inputListener *listener.Listener
	projectionMatrix *mat4.Mat4
}
func Open(w,h int, sWin string, fullscreen bool) *Window {
	
	if GWindow != nil {
		panic("only one Window is allowed")
	}
	listener.GListener = listener.New()
	win := &Window{width: w, height: h, title: sWin, inputListener: listener.GListener}
	GWindow = win
	
	var err error
	if err = glfw.Init(); err != nil {
		panic(fmt.Sprintf("[e] %v\n", err))
	}
	mode := glfw.Windowed
	if fullscreen {mode = glfw.Fullscreen}
	if err = glfw.OpenWindow(win.width, win.height, 8, 8, 8, 0, 16, 0, mode); err != nil {
		panic(fmt.Sprintf("[e] %v\n", err))
	}
	gl.Init() //WHY?? (shader creation fails when not)
	
	glfw.SetSwapInterval(1) // Enable vertical sync on cards that support it.
	glfw.SetWindowTitle(win.title) // Set window title
	//CALLBACKS
	glfw.SetWindowSizeCallback(func (w,h int) {GWindow.reshape(w,h)})
	glfw.SetWindowCloseCallback(listener.OnClose)
	glfw.SetMouseButtonCallback(listener.OnMouseButton)
	glfw.SetMouseWheelCallback (listener.OnMouseWheel)
	glfw.SetKeyCallback (listener.OnKey)
	glfw.SetCharCallback(listener.OnChar)
	
	win.initGL()
	return win
}
func (this *Window) Width () int      {return this.width}
func (this *Window) Height() int      {return this.height}
func (this *Window) ProjectionMatrix() *mat4.Mat4 {return this.projectionMatrix}
func (this *Window) InputListener() *listener.Listener {return this.inputListener}
func (this *Window) Close () {glfw.CloseWindow(); glfw.Terminate()}
func (this *Window) Flip  () {glfw.SwapBuffers()}
func (this *Window) ResetModelViewMatrix() {
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}
func (this *Window) ApplyCamera(cam *camera.Camera) {
	gl.MultMatrixd(cam.Matrix().Ptr())
}
func (this *Window) ClearBuffers () {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT )
}
func (this *Window) initGL() {
	gl.ShadeModel( gl.SMOOTH ) //SMOOTH or FLAT
	gl.CullFace( gl.BACK )
	gl.FrontFace(gl.CCW)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.LIGHTING)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	/*gl.Enable(gl.FOG)
	gl.Fogi(gl.FOG_MODE, gl.EXP)
	gl.Fogfv(gl.FOG_COLOR, []float32{0.5,0.5,0.5,1.0})
	gl.Fogf(gl.FOG_DENSITY, 0.0035)
	gl.Hint(gl.FOG_HINT, gl.DONT_CARE)
	gl.Fogf(gl.FOG_START, 1.0)
	gl.Fogf(gl.FOG_END  , 5000.0)//*/
	gl.Enable(gl.TEXTURE_2D)
	gl.TexEnvf(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.MODULATE) //before: decal, std: modulate
	//gl.TexParameteri(gl.TEXTURE_2D, gl.GENERATE_MIPMAP, true) //mipmaps (dont use it!, its bad!)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR) //GL_LINEAR or GL_NEAREST, no mipmap here
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR) //*/
	gl.ClearColor(.1,.1,.1,1)
	//gl.Enable(gl.PRIMITIVE_RESTART)
	//gl.PrimitiveRestartIndex()
	this.setGLViewport()
}
func (this *Window) setGLViewport() {
	this.reshape(this.width, this.height)
}

func (this *Window) reshape(w,h int) {
	gl.Viewport(0,0,w,h)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	ratio := float64(w)/float64(h)
	znear, zfar, fovx := 0.01, 65536.0, 80.0
	/*var top, bottom, left, right float64
	right = znear * math.Tan(DegToRad(.5*fovx))
	left = -right
	top = right/ratio
	bottom = -top
	gl.Frustum(left,right,bottom,top,znear, zfar)//*/
	this.projectionMatrix = mat4.Projection(DegToRad(fovx), ratio, znear, zfar)
	gl.LoadMatrixd(this.projectionMatrix.Ptr())
	gl.MatrixMode(gl.MODELVIEW)
}



