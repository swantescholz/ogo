package main

import (
	//"github.com/banthar/gl"
	//"github.com/jteeuwen/glfw"
	"fmt"
	//"math"
	//"math/rand"
	"time"
	"io/ioutil"
	"strings"
	"window"
	//"app"
)



func mymain() {
	b,e := ioutil.ReadFile("res/config.txt")
	if e != nil {
		fmt.Println(e)
		return
	}
	var s string = string(b)
	lines := strings.Split(s, "\n")
	var width, height, fullscreen int
	title := lines[2]
	fmt.Sscanf(lines[0], "%v", &width)
	fmt.Sscanf(lines[1], "%v", &height)
	fmt.Sscanf(lines[3], "%v", &fullscreen)
	
	win := window.Open(width,height, title, fullscreen==1)
	defer win.Close()
	win.ClearBuffers()
	win.Flip()
	app := NewApp(win)
	app.Run()
}

func main() {
	fmt.Println("START")
	time1 := time.Now()
	mymain()
	dur := time.Since(time1)
	fmt.Println("END")
	fmt.Println("dt:", dur.Nanoseconds(), "ns", dur.Seconds(), "s")
}


