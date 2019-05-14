package vec3

import (
	"math"
	. "ogo/common"
	. "testing"
)

func TestAngle(t *T) {
	var a,b = V3(0,1,0), V3(1,0,0)
	AssEqual(t, a.Angle(b), Double(math.Pi) / 2)
}

func TestDot(t *T) {
	var a,b = V3(2,3,1), V3(5,3,6)
	AssEqual(t, a.Dot(b), Double(25))
}
