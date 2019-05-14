package common

import (
	. "ogo/globals"
	"math"
	"math/rand"
	"strings"
	. "fmt"
	"os"
)

func Osexit() {
	Println("===---OS_EXIT(1)---===")
	os.Exit(1)
}

type Equaler interface {
	Equal(b Equaler) bool
}

type Byte uint8
type Byte2 uint16
type Byte4 uint32
type Double float64
func (a Byte) Equal(b Equaler) bool {return a == b}
func (a Byte2) Equal(b Equaler) bool {return a == b}
func (a Byte4) Equal(b Equaler) bool {return a == b}
func (a Double) Equal(b Equaler) bool {return a == b}

func (a Double) F    () float64 {return float64(a)}
func (a Double) Set(b Double) {a = b}
func (a Double) SetPrec(b int) {}
func (a Double) SetString(s string) {Sscanf(s,"%v",&a)}
func (a Double) Clear() {}
func (a Double) Add(b,c Double) {a = b + c}
func (a Double) Sub(b,c Double) {a = b - c}
func (a Double) Mul(b,c Double) {a = b * c}
func (a Double) Div(b,c Double) {a = b / c}
func (a Double) Cmp(b Double) int {if a < b {return -1}; if b < a {return 1}; return 0}
func (a Double) Min(b Double) Double {if a < b {return a}; return b}
func (a Double) Max(b Double) Double {if a > b {return a}; return b}
func (a Double) Mod(b Double) Double {return Double(math.Mod(a.F(),b.F()))}
func (a Double) Pow(b Double) Double {return Double(math.Pow(a.F(),b.F()))}
func (a Double) Sqrt () Double {return Double(math.Sqrt(a.F()))}
func (a Double) Floor() Double {return Double(math.Floor(a.F()))}
func (a Double) Log  () Double {return Double(math.Log(a.F()))}
func (a Double) Log2 () Double {return Double(math.Log2(a.F()))}
func (a Double) Log10() Double {return Double(math.Log10(a.F()))}
func (a Double) Sin  () Double {return Double(math.Sin(a.F()))}
func (a Double) Cos  () Double {return Double(math.Cos(a.F()))}
func (a Double) Tan  () Double {return Double(math.Tan(a.F()))}
func (a Double) Asin () Double {return Double(math.Asin(a.F()))}
func (a Double) Acos () Double {return Double(math.Acos(a.F()))}
func (a Double) Atan () Double {return Double(math.Atan(a.F()))}
func (a Double) Atan2(b Double) Double {return Double(math.Atan2(a.F(),b.F()))}
func (a Double) Inv () Double {return Double(1.0 / a)}
func (a Double) Abs () Double {return Double(math.Abs(a.F()))}
func (a Double) AlmostEqual(b Double) bool {return (b-a).Abs() < Epsilon}
func (a Double) ToRad() Double {return a * 0.0174532925199}
func (a Double) ToDeg() Double {return a * 57.295779513082}
func (a Double) Clamp(min,max Double) Double {if a < min {return min}; if a > max {return max}; return a}
func (a Double) Clamp01() Double {return a.Clamp(0.0,1.0)}
func RandDouble() Double {return Double(rand.Float64())}

type String string

func Str(str string) String {return String(str)}
func (s String) S() string {return string(s)}
func (s String) Set(o String) {s = o}
func (s String) Sets(o string) {s = String(o)}
func (s String) Len() int {return len(s.S())}
func (s String) At(index int) String {return Str(s.S()[index:index+1])}
func (s String) Empty() bool {return s.Len() == 0}
func (s String) Substring(a,b int) String {return Str(s.S()[a:b])}
func (s String) SubstringFrom(a int) String {return Str(s.S()[a:])}
func (s String) SubstringTo(b int) String {return Str(s.S()[:b])}
func (s String) BeginsWith(start String) bool {return strings.HasPrefix(s.S(),start.S())}
func (s String) EndsWith(end String) bool {return strings.HasSuffix(s.S(),end.S())}
func (s String) Append(o String) {s = String(s.S() + o.S())}
func (s String) Equal(o Equaler) bool {return s.S() == o.(String).S()}

