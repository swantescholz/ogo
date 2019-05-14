package fractal

import (
	. "github.com/salviati/gmp"
	. "fmt"
	. "ogo/common"
	"ogo/ass"
	. "testing"
)

func iseq(a,b interface{}) bool {return a == b}
type Coord struct {
	a,b int
}

func TestInterfacesBehavesAsExpected(t *T) {
	var a,b Double = 4,4
	ass.True(iseq(a,b))
	var c,d string = "aa","aa"
	ass.True(iseq(c,d))
	var e,f int = 8,8
	ass.True(iseq(e,f))
	var g,h Coord = Coord{2,4},Coord{2,4}
	ass.True(iseq(g,h))
	var i,j *Coord = &g, &h
	ass.False(iseq(i,j))
}

var asses = func(a Stringer, b string) {ass.Equal(Str(a.String()),Str(b))}
var assef = func(x Double, y float64) {ass.Equal(x,Double(y))}
	
func TestGmpString(t *T) {
	var a,b,c,d = NewFloat(50.0),NewFloat(5.0),NewFloat(0.5),NewFloat(0.0625)
	asses(a,"50.0")
	asses(b,"5.0")
	asses(c,"0.5")
	asses(d,"0.0625")
}

func TestGmpStringMinus(t *T) {
	var a,b,c,d = NewFloat(-50.0),NewFloat(-5.0),NewFloat(-0.5),NewFloat(-0.0625)
	asses(a,"-50.0")
	asses(b,"-5.0")
	asses(c,"-0.5")
	asses(d,"-0.0625")
}

func TestDoubleMath(t *T) {
	return
	var a,b,c Double = 1.0,3.0,4	
	a.Mul(b,c)
	assef(a,12)
	c.Add(c,a)
	assef(c,16)
	b.Sub(c,b)
	assef(b,13)
	a.Div(c,a)
	assef(a,.75)
}
