package timefunction

import (
	"fmt"
	"container/list"
	"ogo/mixer/util"
	. "ogo/common"
	"ogo/globals"
)

func Create(term string) (*TimeFunction) {
	tf,e := createTimeFunctionFromString(term)
	if e != nil {
		panic(e)
	}
	return tf
}
func (f *TimeFunction) ToFloats(sps int, length,offset Double) []Double {
	nsamples := int(length*Double(sps))
	data := make([]Double, nsamples)
	for i := range data {
		t := Double(i)*length/Double(nsamples-1) + offset
		y := f.Evaluate(t)
		data[i] = y //before: int16(y*(256*128-8))
	}
	return data
}
func (f *TimeFunction) WriteToWAV(file string, sps int, length,offset Double) {
	data := f.ToFloats(sps, length, offset)
	util.WriteWav(file, data, sps)
}

type ParseError struct {
	Str string
	Err string
}
func (e *ParseError) String() string {
	return fmt.Sprintf("Failed to parse '%v'\n-->ERROR: '%v'", e.Str, e.Err)
}
func (e *ParseError) Error() string {return e.String()}
const (
	opStart = iota //0
	opEnd   = iota //1
	opLB    = iota //2
	opRB    = iota //3
	op0Args = iota //4
	opConst = iota //5
	opTime  = iota //
	opPi    = iota //
	opE     = iota //
	opRand  = iota //
	op1Args = iota //10
	opMir   = iota //
	opMin   = iota //
	opMax   = iota //
	opSqrt  = iota //
	opLn    = iota //15
	opAbs   = iota //
	opFloor = iota //
	opSin   = iota //
	opCos   = iota //
	opTan   = iota //
	op2Args = iota //
	opAdd   = iota //
	opSub   = iota //
	opMul   = iota //
	opDiv   = iota //
	opMod   = iota //
	opPow   = iota //
	opComma = iota //
	//opNArgs = iota
)
type TimeFunction struct {
	childs []*TimeFunction
	//parent *TimeFunction
	op int
	value Double
	isdone bool
}
func (v *TimeFunction) MinChild(t Double) Double {
	if len(v.childs) <= 0 {return 0.0}
	n := len(v.childs[0].childs)
	if n == 0 {
		return 0.0
	}
	m := v.childs[0].childs[0].Evaluate(t)
	for i := 1; i < n; i++ {
		tmp := v.childs[0].childs[i].Evaluate(t)
		if tmp < m {m = tmp}
	}
	return m
}
func (v *TimeFunction) MaxChild(t Double) Double {
	if len(v.childs) <= 0 {return 0.0}
	n := len(v.childs[0].childs)
	if n == 0 {
		return 0.0
	}
	m := v.childs[0].childs[0].Evaluate(t)
	for i := 1; i < n; i++ {
		tmp := v.childs[0].childs[i].Evaluate(t)
		if tmp > m {m = tmp}
	}
	return m
}

func (v *TimeFunction) Evaluate(t Double) Double {
	var a,b Double = 0.0,0.0
	if v.op > op2Args {
		if len(v.childs) < 2 {
			panic(fmt.Sprintf("Too few childs for operator (%v) evaluation (2 needed, %v received", v.op, len(v.childs)))
		} else {
			a = v.childs[0].Evaluate(t)
			b = v.childs[1].Evaluate(t)
		}
	} else if v.op > op1Args {
		if len(v.childs) < 1 {
			panic(fmt.Sprintf("Too few childs for operator (%v) evaluation (1 needed, %v received", v.op, len(v.childs)))
		} else {
			a = v.childs[0].Evaluate(t)
		}
	}
	switch v.op {
		case opConst : return v.value
		case opTime  : return t
		case opPi    : return globals.Pi
		case opE     : return globals.E
		case opRand  : return RandDouble()
		case opMir   : return -a
		case opMin   : return v.MinChild(t)
		case opMax   : return v.MaxChild(t)
		case opSqrt  : return a.Sqrt()
		case opLn    : return a.Log()
		case opAbs   : return a.Abs()
		case opFloor : return a.Floor()
		case opSin   : return a.Sin()
		case opCos   : return a.Cos()
		case opTan   : return a.Tan()
		case opAdd   : return a+b
		case opSub   : return a-b
		case opMul   : return a*b
		case opDiv   : return a/b
		case opMod   : return a.Mod(b)
		case opPow   : return a.Pow(b)
		case opComma : return 0.0 //error
		default : panic(fmt.Sprintf("Operator %v not handled", v.op))
	}
	return 0.0
}
func charIsDigit(c uint8) bool {
	if c == '0' {return true}
	if c == '1' {return true}
	if c == '2' {return true}
	if c == '3' {return true}
	if c == '4' {return true}
	if c == '5' {return true}
	if c == '6' {return true}
	if c == '7' {return true}
	if c == '8' {return true}
	if c == '9' {return true}
	return false
}
func checkStringForNumber(s string) (Double, int) {
	l := len(s)
	if l <= 0 {return 0.0,0}
	if l <= 1 && !charIsDigit(s[0]) {return 0.0,0}
	p,r := -1,l
	for i := 0; i < l; i++ {
		if s[i] == '.' {
			if p < 0 {
				p = i
			} else {
				r = i
				break
			}
		} else if !charIsDigit(s[i]) {
			r = i
			break
		}
	}
	sub := s[:r]
	var f Double = 0.0
	fmt.Sscanf(sub, "%f", &f)
	return f, len(sub)
}
func createTimeFunctionFromString(term string) (tf *TimeFunction, err error) {
	term = term + "##########" //ending chars
	ops  := make([]int, 0, len(term))
	nums := make([]Double, 0, len(term))
	nLB := 0
	//checking braces
	for _,v := range term {
		if v == '(' {
			nLB++
		} else if v == ')' {
			nLB--
		}
		if nLB < 0 {
			err = &ParseError{term, "Bad parentheses (RB)"}
			return
		}
	}
	if nLB > 0 {
		err = &ParseError{term, "Bad parentheses (LB)"}
		return
	}
	//scanning
	for i := 0; i < len(term); i++ {
		if term[i] == '#' {break}
		if        term[i:i+5] == "floor" {ops = append(ops, opFloor); i += 5-1; continue
		} else if term[i:i+4] == "sqrt"  {ops = append(ops, opSqrt ); i += 4-1; continue
		} else if term[i:i+3] == "sin"   {ops = append(ops, opSin  ); i += 3-1; continue
		} else if term[i:i+3] == "cos"   {ops = append(ops, opCos  ); i += 3-1; continue
		} else if term[i:i+3] == "tan"   {ops = append(ops, opTan  ); i += 3-1; continue
		} else if term[i:i+3] == "abs"   {ops = append(ops, opAbs  ); i += 3-1; continue
		} else if term[i:i+3] == "min"   {ops = append(ops, opMin  ); i += 3-1; continue
		} else if term[i:i+3] == "max"   {ops = append(ops, opMax  ); i += 3-1; continue
		} else if term[i:i+2] == "ln"    {ops = append(ops, opLn   ); i += 2-1; continue
		} else {
			f,n := checkStringForNumber(term[i:])
			if n > 0 {
				ops = append(ops, opConst)
				nums = append(nums, f)
				i += n-1
				continue
			}
		}
		switch term[i] {
			case '#': break // # is ending char
			case '(': ops = append(ops, opLB)
			case ')': ops = append(ops, opRB)
			case 't': ops = append(ops, opTime)
			case 'p': ops = append(ops, opPi)
			case 'e': ops = append(ops, opE)
			case 'r': ops = append(ops, opRand)
			case '+': ops = append(ops, opAdd)
			case '-': ops = append(ops, opSub)
			case '*': ops = append(ops, opMul)
			case '/': ops = append(ops, opDiv)
			case '%': ops = append(ops, opMod)
			case '^': ops = append(ops, opPow)
			case ',': ops = append(ops, opComma)
			default:
		}
	}
	//(-xy) handling
	for i,v := range ops {
		if v == opSub {
			if i == 0 {
				ops[i] = opMir
			} else if ops[i-1] == opLB {
				ops[i] = opMir
			}
		}
	}
	funcs := make([]*TimeFunction, 0, len(ops))
	for _,v := range ops {
		var f Double = 0.0
		if v == opConst {
			if len(nums) <= 0 {panic("Too few nums!")}
			f = nums[0]
			if len(nums) >= 2 {nums = nums[1:]}
		}
		tmp := &TimeFunction{op: v, value: f}
		switch tmp.op {
			case opConst: fallthrough
			case opTime : fallthrough
			case opPi   : fallthrough
			case opE    : fallthrough
			case opRand :
				tmp.isdone = true
		}
		funcs = append(funcs, tmp)
	}
	l := list.New()
	l.PushBack(&TimeFunction{op: opStart})
	for _,v := range funcs {
		l.PushBack(v)
	}
	l.PushBack(&TimeFunction{op: opEnd})
	tf = simplifyTimeFunctionTokens(l)
	return
}
func simplifyBracelessTimeFunctionTokens(funcs *list.List, a *list.Element, b *list.Element) {
	for a.Next().Next() != b {
		for i := 0; i <= 6; i++ {
			for e := a; e != b; e = e.Next() {
				v := e.Value.(*TimeFunction)
				if v.op == opStart || v.op == opEnd {continue}
				x := e.Prev().Value.(*TimeFunction)
				y := e.Next().Value.(*TimeFunction)
				switch i {
				case 0:
					if v.op == opMir {
						if y.isdone {
							v.childs = []*TimeFunction{y}
							funcs.Remove(e.Next())
							v.isdone = true
							i = -1; break
						}
					}
				case 1:
					if v.op == opMin || v.op == opMax {
						if y.isdone {
							v.childs = []*TimeFunction{y}
							funcs.Remove(e.Next())
							v.isdone = true
							i = -1; break
						}
					}
				case 2:
					if v.op == opAbs || v.op == opFloor || v.op == opSqrt || v.op == opLn ||
						v.op == opSin || v.op == opCos || v.op == opTan {
						if y.isdone {
							v.childs = []*TimeFunction{y}
							funcs.Remove(e.Next())
							v.isdone = true
							i = -1; break
						}
					}
				case 3:
					if v.op == opPow {
						if x.isdone && y.isdone {
							v.childs = []*TimeFunction{x,y}
							funcs.Remove(e.Prev())
							funcs.Remove(e.Next())
							v.isdone = true
							i = -1; break
						}
					}
				case 4:
					if v.op == opMul || v.op == opDiv || v.op == opMod {
						if x.isdone && y.isdone {
							v.childs = []*TimeFunction{x,y}
							funcs.Remove(e.Prev())
							funcs.Remove(e.Next())
							v.isdone = true
							i = -1; break
						}
					}
				case 5:
					if v.op == opAdd || v.op == opSub {
						if x.isdone && y.isdone {
							v.childs = []*TimeFunction{x,y}
							funcs.Remove(e.Prev())
							funcs.Remove(e.Next())
							v.isdone = true
							i = -1; break
						}
					}
				case 6:
					if v.op == opComma {
						if x.isdone && y.isdone {
							v.childs = []*TimeFunction{x,y}
							funcs.Remove(e.Prev())
							funcs.Remove(e.Next())
							v.isdone = true
							i = -1; break
						}
					}
				}
			}
		}
	}
}
func simplifyTimeFunctionTokens(funcs *list.List) *TimeFunction {
	var lb, rb *list.Element
	for {
		lb, rb = nil,nil
		for e := funcs.Front(); e != nil; e = e.Next() {
			if e.Value.(*TimeFunction).op == opLB {lb = e}
			if e.Value.(*TimeFunction).op == opRB {
				rb = e
				simplifyBracelessTimeFunctionTokens(funcs, lb,rb)
				funcs.Remove(lb)
				funcs.Remove(rb)
				break
			}
		}
		if rb == nil && lb == nil {
			if funcs.Len() <= 0 {panic("Too few funcs")}
			simplifyBracelessTimeFunctionTokens(funcs, funcs.Front(), funcs.Back())
			return funcs.Front().Next().Value.(*TimeFunction)
		}
	}
	panic("unexpected situation")
}


