package ass

import (
	. "fmt"
	. "ogo/globals"
	. "ogo/common"
	"ogo/util"
	"runtime"
	"strings"
	"log"
)

func Stacktrace() (trace string, fatalFile string, fatalLine int) {
	trace, fatalLine = "", -1
	for skip := MinStacktraceDepth; skip <= StacktraceSize + MinStacktraceDepth; skip++ {
		_,file,line,ok := runtime.Caller(skip)
		if !ok {break}
		if fatalLine < 0 {
			fatalFile = file
			fatalLine = line
		}
		file = removePathFromFileIfLocal(file)
		trace += Sprintf("\nline %v in %v", util.IndentNumber(line, 4), file)
	}
	trace = trace[1:]
	return 
}

func printStacktrace(trace string) {
	Println("----- STACKTRACE -----")
	Println(trace)
	Println("----------------------")
}

func removePathFromFileIfLocal(filepath string) string {
	if !strings.HasPrefix(filepath, "/usr") {
		return util.RemovePathFromFile(filepath)
	}
	return filepath
}

func assMsg(file string, line int, a,b interface{}, op string) string {
	var head = "====-ASSERTION-FAILED-===="
	var format = "\n%v\n%v\nline %v:%v\n%v %v %v"
	var subs = strings.Repeat("-",len(head))
	var code = Nl + subs + Nl +
		strings.Trim(util.ReadCodeLineFromFile(file, line), " \t") +
		Nl + subs
	return Sprintf(format, head,removePathFromFileIfLocal(file),line,code,a,op,b)
}

func assFatal(a,b interface{}, op string) {
	trace, fatalfile, fatalline := Stacktrace()
	printStacktrace(trace)
	log.Fatal(assMsg(fatalfile,fatalline,a,b,op))
}

func Equal(a,b Equaler) {
	if !a.Equal(b) {
		assFatal(a,b, "!=")
	}
}
func Unequal(a,b Equaler) {
	if a.Equal(b) {
		assFatal(a,b, "==")
	}
}
func True(got bool) {
	if !got {
		assFatal(got, "true", "!=")
	}
}
func False(got bool) {
	if got {
		assFatal(got, "false", "!=")
	}
}
func Error(got error) {
	if got != nil {
		assFatal(got, "nil", "!=")
	}
}

