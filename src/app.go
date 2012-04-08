package main

import (
	"github.com/banthar/gl"
	"github.com/jteeuwen/glfw"
	"fmt"
	//"math"
	//"math/rand"
	"time"
	"window"
	"input/listener"
	. "input/constants"
	"camera"
	"glmath/mat4"
	. "glmath/vec3"
	. "glmath/util"
	. "glmath/color"
	"glmath/material"
	"texture"
	"program"
	"model"
)

var T0 float64 = 0
var Frames uint32 = 0

type App struct {
	win *window.Window
	il *listener.Listener
}

func NewApp(w *window.Window) *App {
	return &App{w, w.InputListener()}
}

func psleep(s interface{}) {
	fmt.Printf("%v\n", s)
	dur,_	 := time.ParseDuration("0.3s")
	time.Sleep(dur)
}

func (this *App) Run() {
	fmt.Println("running")

	var done bool
	cam := camera.New()
	cam.Pos = V3(-0,0,3)
	cam.Dir = V3(0,-0,-1)
	
	pos := []float32{5.0, 5.0, 10.0, 0.0}
	gl.Lightfv(gl.LIGHT0, gl.POSITION, pos)
	gl.Enable(gl.LIGHT0)
	
	elapsed := 0.0
	rotateSpeed := DegToRad(100.0)
	moveSpeed   := 3.0
	texture.SetDefaultTexturePaths("res/textures/", ".png")
	texture.Load("ground")
	texture.Load("ground2")
	defer texture.ClearAllData()
	tex1 := texture.Get("ground2")
	//tex2 := texture.Get("ground2")
	model.SetDefaultModelPaths("res/models/", ".obj")
	model.Load("teapot")
	model.Load("cuboid")
	model.Load("sphere")
	model.Load("bunny")
	defer model.ClearAllData()
	mod1 := model.Get("bunny")//*/
	
	//*
	prog1 := program.New("test", "test")
	defer prog1.Destroy()
	prog := prog1
	println(prog.InfoLog())
	var mView, mProjection, mModel *mat4.Mat4
	mModel = mat4.New()//*/
	
	done = false
	for !done {
		time1 := time.Now()
		this.win.ClearBuffers()
		this.win.ResetModelViewMatrix()
		
		//input
		if this.il.KeyDown(glfw.KeyLeft ) {cam.RotateY(-rotateSpeed * elapsed)}
		if this.il.KeyDown(glfw.KeyRight) {cam.RotateY( rotateSpeed * elapsed)}
		if this.il.KeyDown(glfw.KeyUp  ) {cam.RotateX (-rotateSpeed * elapsed)}
		if this.il.KeyDown(glfw.KeyDown) {cam.RotateX ( rotateSpeed * elapsed)}
		if this.il.KeyDown(KeyX) {cam.MoveZ(-moveSpeed * elapsed)}
		if this.il.KeyDown(KeyC) {cam.MoveZ( moveSpeed * elapsed)}
		if this.il.KeyPressed(KeyR) {prog.Reload(); println(prog.InfoLog())}
		
		this.win.ApplyCamera(cam)
		//*
		mView = cam.Matrix()
		mProjection = this.win.ProjectionMatrix()
		mModel.SetIdentity()
		prog.SetModelMatrix     (mModel)
		prog.SetViewMatrix      (mView)
		prog.SetProjectionMatrix(mProjection)
		prog.UniformColor("uColor", Col(1,.1,0.5,1))
		prog.UniformMat4("uMat", mModel)
		prog.Use()//*/
		material.White.Use()
		//tex1.Unbind()
		mod1.Render()
		tex1.Bind2(Col(1,1,1,1))
		gl.Begin(gl.TRIANGLES)
		gl.Normal3d(0,0,1)
		gl.TexCoord2d(11.0, 10.0)
		gl.Vertex3d(1.0,0.0 ,-15.0)
		gl.TexCoord2d(0.5, 11.0)
		gl.Vertex3d(0.0,1.0 ,-15.0)
		gl.TexCoord2d(0.0, 0.0)
		gl.Vertex3d(-1.0,0.0,-15.0)
		gl.TexCoord2d(1.0, 0.0)
		gl.Vertex3d(1.0,0.0,-5.0)
		gl.TexCoord2d(0.5, 1.0)
		gl.Vertex3d(0.0,1.0,-5.0)
		gl.TexCoord2d(0.0, 0.0)
		gl.Vertex3d(-1.0,0.0,-5.0)
		gl.End()//*/
		this.win.Flip()
		done = this.il.MBQuit() || this.il.KeyPressed(glfw.KeyEsc)
		elapsed = time.Since(time1).Seconds()
		
		printFPS()
		
	}
}



//############################
func printFPS() {
	Frames++
	{
		t := glfw.Time()
		if t-T0 >= 5 {
			seconds := (t - T0)
			fps := float64(Frames) / seconds
			print(Frames, " frames in ", int(seconds), " seconds = ", int(fps), " FPS\n")
			T0 = t
			Frames = 0
		}
	}
}
	

