package program

import (
	"github.com/banthar/gl"
	"fmt"
	"strings"
	"io/ioutil"
)


var (
	GDefaultShaderPathStart     = "res/shader/"
	GDefaultVertShaderPathStart = ""
	GDefaultFragShaderPathStart = ""
	GDefaultVertShaderPathEnd   = ".vert"
	GDefaultFragShaderPathEnd   = ".frag"
)
func SetDefaultShaderSearchPaths(s0, sv, sf, ev, ef string) {
	GDefaultShaderPathStart     = s0
	GDefaultVertShaderPathStart = sv
	GDefaultFragShaderPathStart = sf
	GDefaultVertShaderPathEnd   = ev
	GDefaultFragShaderPathEnd   = ef
}

type Shader struct {
	id gl.Shader
	isVert bool
}
func NewShader(isvert bool) *Shader {
	type_ := gl.GLenum(gl.FRAGMENT_SHADER)
	if isvert {type_ = gl.GLenum(gl.VERTEX_SHADER)}
	id := gl.CreateShader(type_)
	sh := &Shader{id, isvert}
	return sh
}
func (s *Shader) Destroy() {s.id.Delete()}
func (s *Shader) IsVertShader() bool {return  s.isVert}
func (s *Shader) IsFragShader() bool {return !s.isVert}
func (s *Shader) InfoLog() string {return s.id.GetInfoLog()}
func (s *Shader) GetSource () string {return s.id.GetSource ()}
func (s *Shader) Source    (source string) {s.id.Source(source)}
func (s *Shader) Compile() {s.id.Compile()}
func (s *Shader) String() string {
	return fmt.Sprintf("%v", s.id)
}

func stringEndsWith(s, end string) bool {
	if len(end) > len(s) {return false}
	return s[len(s)-len(end):] == end
}
func ExtendVertShaderSourcePath(srcpath string) string {
	return GDefaultShaderPathStart + GDefaultVertShaderPathStart +
		srcpath + GDefaultVertShaderPathEnd
}
func ExtendFragShaderSourcePath(srcpath string) string {
	return GDefaultShaderPathStart + GDefaultFragShaderPathStart +
		srcpath + GDefaultFragShaderPathEnd
}

func PreprocessShaderSource(src string) string {
	const sinclude = "//#include "
	for {
		pos := strings.Index(src, sinclude)
		if pos < 0 {break}
		pos += len(sinclude)
		inewline := strings.Index(src[pos:], "\n")+pos
		var path = ""
		if inewline < 0 {
			path = src[pos:]
		} else {
			path = src[pos: inewline]
		}
		path = GDefaultShaderPathStart + strings.TrimSpace(path)
		pos -= len(sinclude)
		b,e := ioutil.ReadFile(path)
		if e != nil {
			panic(e)
		}
		extrasrc := string(b)
		endsrc := ""
		if inewline >= 0 { endsrc = src[inewline:]}
		src = src[:pos] + extrasrc + endsrc
		//println(src)
	}
	return src
}

func vertShaderSource(sh string) string {
	if !strings.Contains(sh, "\n") { // && !strings.Contains(sh, "main") {
		sh = ExtendVertShaderSourcePath(sh)
		b,e := ioutil.ReadFile(sh)
		if e != nil {
			panic(e)
		}
		return string(b)
	}
	return sh
}
func fragShaderSource(sh string) string {
	if !strings.Contains(sh, "\n") { // && !strings.Contains(sh, "main") {
		sh = ExtendFragShaderSourcePath(sh)
		b,e := ioutil.ReadFile(sh)
		if e != nil {
			panic(e)
		}
		return string(b)
	}
	return sh
}
