package model

import (
	"github.com/banthar/gl"
	. "fmt"
	"ogo/glmath/mat4"
	. "ogo/glmath/vec3"
	"ogo/program"
	"ogo/globals"
	"math"
	
)

type Model struct {
	data *ModelData
	Mat *mat4.Mat4
}

func NewModel(md *ModelData) *Model {
	return &Model{data: md, Mat: mat4.New()}
}

func (m *Model) Data () *ModelData {return m.data}
func (m *Model) Render() {
	if globals.UseShader {
		program.SetModelMatrix(m.Mat)
		m.data.Render()
	} else {
		gl.PushMatrix()
		gl.MultMatrixf(m.Mat.Ptr32())
		m.data.Render()
		gl.PopMatrix()
	}
}

func (m *Model) String() string {
	return Sprintf("%v", m.data.String())
}

func RenderArrow(pos, dir *Vec3) {
	const tolerance = 0.0001
	arrow := Get("arrow")
	l := dir.Length()
	d0 := dir.Unit()
	zaxis := V3(0,0,1)
	arrow.Mat = mat4.Scaling(V3(.1*l,.1*l,.5*l))
	angle := zaxis.Angle(d0)
	if angle > tolerance && math.Pi-angle > tolerance {
		arrow.Mat = arrow.Mat.Mul(mat4.RotationAxis(zaxis.Add(d0).Unit(), math.Pi))
	} else if math.Pi-angle <= tolerance {
		arrow.Mat = arrow.Mat.Mul(mat4.Scaling(V3(1,1,-1)))
	}
	arrow.Mat = arrow.Mat.Mul(mat4.Translation(pos.Add(dir.Muls(.5))))
	arrow.Render()
}

