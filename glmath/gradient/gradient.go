package gradient

import (
	. "fmt"
	. "ogo/glmath/color3"
	. "ogo/common"
	. "ogo/globals"
	"ogo/util"
)

type Gradient struct {
	colors []*Color
}

func NewGradient(colors []*Color) *Gradient {
	return &Gradient{colors: colors}
}

func (this *Gradient) String() string {
	return Sprintf("%v", this.Length())
}
func (this *Gradient) Length() int {return len(this.colors)}

func (this *Gradient) findIntervalIndexOfFactor(factor Double) int {
	return int(factor * (Double(this.Length()) - 1.0 - Epsilon))
}

func (this *Gradient) projectFactorIntoInterval(factor Double, intervalindex int) Double {
	return factor * (Double(this.Length()) - 1.0) - Double(intervalindex)
}

func (this *Gradient) Interpolate(factor Double) *Color {
	factor = util.Clamp01(factor)
	var index = this.findIntervalIndexOfFactor(factor)
	var newfactor = this.projectFactorIntoInterval(factor, index)
	return this.colors[index].Interpolated(this.colors[index+1], newfactor)
}

