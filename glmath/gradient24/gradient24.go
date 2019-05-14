package gradient

import (
	. "fmt"
	. "ogo/glmath/color24"
	"ogo/glmath/color"
	. "ogo/common"
	. "ogo/globals"
	"ogo/util"
)

type Gradient struct {
	colors []*Color24
	length int
	elength Double
}

func NewGradient(colors []*color.Color) *Gradient {
	var this = &Gradient{}
	var l = len(colors)
	this.colors = make([]*Color24, l)
	for i := 0; i < l; i++ {
		this.colors[i] = colors[i].Color24()
	}
	this.length = l
	this.elength = (Double(this.length) - 1.0 - Epsilon)
	return this
}

func (this *Gradient) String() string {
	return Sprintf("%v", this.length)
}
func (this *Gradient) Length() int {return this.length}

func (this *Gradient) findIntervalIndexOfFactor(factor Double) int {
	return int(factor * this.elength)
}

func (this *Gradient) projectFactorIntoInterval(factor Double, intervalindex int) Double {
	return factor * this.elength - Double(intervalindex)
}

var index int
var newfactor Double
var a,b *Color24

func (this *Gradient) Interpolate(col *Color24, factor Double) {
	factor = util.Clamp01(factor)
	index = this.findIntervalIndexOfFactor(factor)
	newfactor = this.projectFactorIntoInterval(factor, index)
	a,b = this.colors[index], this.colors[index+1]
	col.R = uint8(Double(a.R) + (Double(b.R)-Double(a.R))*newfactor)
	col.G = uint8(Double(a.G) + (Double(b.G)-Double(a.G))*newfactor)
	col.B = uint8(Double(a.B) + (Double(b.B)-Double(a.B))*newfactor)
}

