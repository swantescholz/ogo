package model

import (
	//"github.com/banthar/gl"
	"fmt"
	//"math"
	//. "glmath/color"
	//"glmath/mat4"
	//. "glmath/util"
)

type Model struct {
	data *ModelData
	//mat mat4.Mat4
}

func NewModel(md *ModelData) *Model {
	return &Model{data: md}
}

func (m *Model) Data () *ModelData {return m.data}
func (m *Model) Render() {
	m.data.Render()
}

func (m *Model) String() string {
	return fmt.Sprintf("%v", m.data.String())
}



