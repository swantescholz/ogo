package color32

type Color32 struct {
	R,G,B,A uint8
}

func NewColor32(r,g,b,a uint8) *Color32 {
	return &Color32{r,g,b,a}
}
