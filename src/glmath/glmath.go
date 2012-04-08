package glmath

import (
	"github.com/banthar/gl"
	"glmath/color"
	//"glmath/mat3"
	"glmath/mat4"
	//"glmath/qtrnn"
	//"glmath/util"
	"glmath/vec2"
	"glmath/vec3"
)

func GlColor    (c *color.Color) {gl.Color4d(c.R,c.G,c.B,c.A)}
func GlLoadMat4 (m *mat4.Mat4  ) {gl.LoadMatrixd(m.Ptr())}
func GlMultMat4 (m *mat4.Mat4  ) {gl.MultMatrixd(m.Ptr())}
func GlTex2     (v *vec2.Vec2  ) {gl.TexCoord2d(v.X, v.Y)}
func GlVec3     (v *vec3.Vec3  ) {gl.Vertex3d  (v.X, v.Y, v.Z)}
