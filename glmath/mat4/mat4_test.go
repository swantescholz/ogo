package mat4

import (
	. "ogo/common"
	. "testing"
	"ogo/ass"
)

func TestProjection(t *T) {
	var ratio Double = 1.2
	var znear, zfar, fovx Double = 0.01, 65536.0, 80.0
	ass.True(fovx.ToRad() < 7.0)
	var m = Projection(fovx.ToRad(), ratio, znear, zfar)
	t.Logf("projection works: %v\n",m)
}
