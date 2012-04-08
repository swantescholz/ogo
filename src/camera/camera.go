package camera

import (
	"fmt"
	//"math"
	. "glmath/vec3"
	"glmath/mat4"
)

type Camera struct {
	Pos, Dir, Up *Vec3
}

func New() *Camera {
	return &Camera{V3(0,0,0), V3(0,0,-1), V3(0,1,0)}
}

func (this *Camera) MoveZ(d float64) {
	this.Pos = this.Pos.Add(this.Dir.Unit().Muls(d))
}
func (this *Camera) RotateX(angle float64) {
	xaxis := this.Dir.Cross(this.Up)
	mrot := mat4.RotationAxis(xaxis, angle)
	this.Dir = mrot.Mulv(this.Dir)
}
func (this *Camera) RotateY(angle float64) {
	yaxis := this.Up
	mrot := mat4.RotationAxis(yaxis, angle)
	this.Dir = mrot.Mulv(this.Dir)
}

func (this *Camera) Matrix() *mat4.Mat4 {
	return mat4.Camera(this.Pos, this.Dir, this.Up)
}

func (this *Camera) String() string {
	return fmt.Sprintf("")
}
