package util

import (
	"github.com/banthar/gl"
	"image"
	"image/png"
	"image/color"
	. "fmt"
	"strings"
	"os"
	"io/ioutil"
	"time"
	. "ogo/common"
)


//----------------------------------------------------------------------------
func GlErrorString() string {
	code := gl.GetError()
	if code == 0 {
		return "NO GL ERROR"
	}
	index := uint32(code) - uint32(gl.INVALID_ENUM)
	errorStrings := []string{
		"GL_INVALID_ENUM",
		"GL_INVALID_VALUE",
    	"GL_INVALID_OPERATION",
    	"GL_STACK_OVERFLOW",
    	"GL_STACK_UNDERFLOW",
    	"GL_OUT_OF_MEMORY",
    	"GL_UNKNOWN_ERROR_CODE"}
	if (0 <= index && index < 6) {
		return errorStrings[index]
	}
	return errorStrings[6]
}

func GlVersion() (major,minor int) {
	var version = gl.GetString(gl.VERSION)
	Sscan(version, &major)
	Sscan(version[strings.Index(version,".")+1:], &minor)
	return
}

func GlSupportsVersion(minmajor, minminor int) bool {
	var major,minor = GlVersion()
	if major < minmajor {return false}
	if major > minmajor {return true}
	return minor >= minminor
}

func Ps(s interface{}) {
	Printf("%v\n", s)
	dur,_	 := time.ParseDuration("0.3s")
	time.Sleep(dur)
}

var __tsLastTime time.Time
func init() {
	__tsLastTime = time.Now()
}
func Ts() {
	var now = time.Now()
	Println(now.Sub(__tsLastTime))
	__tsLastTime = now
}

func Clamp(f, min, max Double) Double {
	if f < min {return min}
	if f > max {return max}
	return f
}
func Clamp01(f Double) Double {
	if f < 0.0 {return 0.0}
	if f > 1.0 {return 1.0}
	return f
}

func IsPot(n int) bool {return (n & (n-1) == 0)}
func NextLowerPot(n int) int {
	for i := uint32(1); i < 31; i++ {
		if 1<<i > n {return 1<<(i-1)}
	}
	return 1<<30
}

func Wait(sec float64) {
	s := Sprintf("%vs", sec)
	dur,_ := time.ParseDuration(s)
	time.Sleep(dur)
}

func NowAsString() string {
	var now = time.Now()
	var year,month,day = now.Date()
	var hour,minute,second = now.Hour(), now.Minute(), now.Second()
	return Sprintf("%04d-%02d-%02d_%02d-%02d-%02d",year,int(month),day,hour,minute,second)
}

func PixelDataToImage(data []uint8, w,h int) image.Image { 
	if w*h*4 != len(data) {panic("bad data length")}
	
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			yy := h - y - 1
			r := data[4*(yy*w+x)+0]
			g := data[4*(yy*w+x)+1]
			b := data[4*(yy*w+x)+2]
			a := data[4*(yy*w+x)+3]
			rgba.Set(x,y, color.RGBA{r,g,b,a})
		}
	}
	return rgba
}

func WritePng(pngImage image.Image, path string) {
	if !strings.HasSuffix(path, ".png") {path += ".png"}
	f, err := os.OpenFile(path, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {panic(err)}
	defer f.Close()
	err = png.Encode(f, pngImage)
	if err != nil {panic(err)}
}

func IndentNumber(number, width int) string {
	var str = Sprint(number)
	if times := width - len(str); times > 0 {
		str = strings.Repeat(" ", times) + str
	}
	return str
}

func RemovePathFromFile(filepath string) string {
	var i = strings.LastIndex(filepath, "/")
	return filepath[i+1:]
}



func ReadCodeLineFromFile(filepath string, line int) string {
	var content,err = ioutil.ReadFile(filepath)
	if err != nil {panic(err)}
	var lines = strings.Split(string(content), "\n")
	return lines[line-1]
}

