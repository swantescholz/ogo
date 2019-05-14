package camera

import (
	"github.com/banthar/gl"
	. "fmt"
	. "ogo/glmath/vec3"
	"ogo/glmath/mat4"
	"ogo/globals"
	"ogo/program"
	"math"
	. "ogo/common"
)

type Camera struct {
	Pos, Dir, Up *Vec3
}

func (this *Camera) String() string {
	return Sprintf("")
}

//satisfy program.ICamera
func (t *Camera) Position () *Vec3 {return t.Pos}
func (t *Camera) Direction() *Vec3 {return t.Dir}
func (t *Camera) UpVector () *Vec3 {return t.Up}

func New() *Camera {
	return &Camera{V3(0,0,0), V3(0,0,-1), V3(0,1,0)}
}

//BEFORE program.Use()!
func (this *Camera) Use() {
	if globals.UseShader {
		program.SetCamera(this)
	}
	gl.MultMatrixd(this.Matrix().Ptr())
}

func (this *Camera) MoveZ(d Double) {
	this.Pos = this.Pos.Add(this.Dir.Unit().Muls(d))
}
func (this *Camera) RotateX(angle Double) {
	const limitAngle = 0.05
	upAngle := this.Up.Angle(this.Dir)
	if upAngle < limitAngle && angle < 0.0 {
		return
	}
	if (math.Pi - upAngle) < limitAngle && angle > 0.0 {
		return
	}
	xaxis := this.Dir.Cross(this.Up)
	mrot := mat4.RotationAxis(xaxis, angle)
	this.Dir = mrot.Mulv(this.Dir)
}
func (this *Camera) RotateY(angle Double) {
	yaxis := this.Up
	mrot := mat4.RotationAxis(yaxis, angle)
	this.Dir = mrot.Mulv(this.Dir)
}

func (this *Camera) Matrix() *mat4.Mat4 {
	return mat4.Camera(this.Pos, this.Dir, this.Up)
}




