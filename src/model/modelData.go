package model

import (
	"github.com/banthar/gl"
	"fmt"
	//"math"
	"strings"
	//. "glmath/color"
	"glmath/mat4"
	. "glmath/vec2"
	. "glmath/vec3"
	//. "glmath/util"
	"vbo"
	"io/ioutil"
)


type ModelData struct {
	path string
	vbo *vbo.Vbo
	stdMat *mat4.Mat4
	numv,numvn,numvt,numf,numi int32
}

func NewModelData(path string) *ModelData {
	md := &ModelData{path: path, vbo: nil}
	md.stdMat = mat4.New()
	md.vbo = vbo.New()
	md.LoadFile()
	return md
}

//clears the opengl texture memory
func (m *ModelData) Destroy() {
	m.vbo.Destroy()
}

func (m *ModelData) Name  () string {return m.path}
func (m *ModelData) SetStdMatrix(mat *mat4.Mat4) {m.stdMat.Set(mat)}
func (m *ModelData) Render() {
	m.vbo.Enable()
	m.vbo.Draw(gl.TRIANGLES)
	m.vbo.Disable()
}
func (m *ModelData) LoadFile() {
	fileend := m.path[strings.LastIndex(m.path, ".")+1:]
	switch fileend {
		case "obj": m.LoadObj()
		default: panic("model file '."+fileend+"' format not implemented!")
	}
}
func extractNumsFromFaceString(str string) []int32 {
	str = strings.TrimSpace(str)
	parts := strings.Split(str, " ")
	nums := make([]int32, len(parts))
	for i := range parts {
		fmt.Sscan(parts[i], &nums[i])
	}
	return nums
}
func (m *ModelData) LoadObj() {
	by,e := ioutil.ReadFile(m.path)
	if e != nil {
		panic(e)
		return
	}
	var s string = string(by)
	lines0 := strings.Split(s, "\n")
	lines := make([]string, 0, len(lines0))
	m.numv, m.numvn, m.numvt, m.numf, m.numi = 0,0,0,0,0
	for i := range lines0 {
		line := strings.TrimSpace(lines0[i])
		line = strings.Replace(line, "//", " ", -1)
		line = strings.Replace(line, "/", " ", -1)
		if len(line) < 3 {continue}
		goodLine := true
		switch line[:2] {
			case "v ": m.numv++
			case "vn": m.numvn++
			case "vt": m.numvt++
			case "f ": m.numf++; m.numi += int32((strings.Count(line, " ")-2)*3)
			default: goodLine = false
		}
		if goodLine {
			lines = append(lines, line)
		}
	}
	usen := m.numvn > 0
	uset := m.numvt > 0
	numvertexparts := 1
	if usen {numvertexparts++}
	if uset {numvertexparts++}
	p := make([]*Vec3, 0, m.numv)
	n := make([]*Vec3, 0, m.numvn)
	t := make([]*Vec2, 0, m.numvt)
	m.vbo.P = make([]*Vec3, 0, m.numi)
	if usen {m.vbo.N = make([]*Vec3, 0, m.numi)}
	if uset {m.vbo.T = make([]*Vec2, 0, m.numi)}
	var x,y,z float64
	for i, line := range lines {
		switch line[:2] {
		case "v ":
			fmt.Sscan(line[2:], &x,&y,&z)
			p = append(p, V3(x,y,z))
		case "vn":
			fmt.Sscan(line[3:], &x,&y,&z)
			n = append(n, V3(x,y,z))
		case "vt":
			fmt.Sscan(line[3:], &x,&y)
			t = append(t, V2(x,y))
		case "f ":
			nums := extractNumsFromFaceString(lines[i][2:])
			lastBackfall := 0
			for j := 0; j < len(nums); {
				m.vbo.P = append(m.vbo.P, p[nums[j+0]-1])
				if uset {
					m.vbo.T = append(m.vbo.T, t[nums[j+1]-1])
					if usen {
						m.vbo.N = append(m.vbo.N, n[nums[j+2]-1])
					}
				} else if usen {
					m.vbo.N = append(m.vbo.N, n[nums[j+1]-1])
				}
				//triangularizing polygons:
				j += numvertexparts
				if j >= len(nums) {break}
				if j == (lastBackfall+3)*numvertexparts {
					lastBackfall++
					j = lastBackfall*numvertexparts
				}
			}
		}
	}
	m.vbo.Create(gl.STATIC_DRAW)
}

func (m *ModelData) String() string {
	return fmt.Sprintf("%v, %v", m.path, m.vbo.String())
}



