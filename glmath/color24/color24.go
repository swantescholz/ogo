package color24

import (
	"github.com/banthar/gl"
)

type Color24 struct {
	R,G,B uint8
}

func NewColor24(r,g,b uint8) *Color24 {
	return &Color24{r,g,b}
}

func (this *Color24) Gl() {
	gl.Color3ub(this.R,this.G,this.B)
}

var (
	White  = &Color24{255,255,255}
	Grey   = &Color24{128,128,128}
	Black  = &Color24{0,0,0}
	Red    = &Color24{255,0,0}
	Green  = &Color24{0,255,0}
	Blue   = &Color24{0,0,255}
	Yellow = &Color24{255,255,0}
	Purple = &Color24{255,0,255}
	Cyan   = &Color24{0,255,255}
	Orange = &Color24{255,128,0}
	Pink   = &Color24{255,51,204}
	SkyBlue      = &Color24{128,174,220}
	MellowOrange = &Color24{220,174,128}
	ForestGreen  = &Color24{64,128,85}
	Silver = &Color24{200,200,200}
	Gold   = &Color24{210,190,0}
)
