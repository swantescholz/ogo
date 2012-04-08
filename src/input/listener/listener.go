package listener

import (
	"github.com/jteeuwen/glfw"
	"fmt"
)

const (
	keyUp = iota
	keyPressed = iota
	keyDown = iota
)

var GListener *Listener = nil

type Listener struct {
	quit bool
	mKeys map[int]int
	mButtons map[int]int
}
func New() *Listener {
	il := &Listener{false, make(map[int]int), make(map[int]int)}
	return il
}

func (this *Listener) MBQuit() bool {
	return this.quit || glfw.WindowParam(glfw.Opened) != 1
}

func (this *Listener) KeyDown(key int) bool {
	return glfw.Key(key) == glfw.KeyPress
}
func (this *Listener) KeyUp(key int) bool {
	return glfw.Key(key) == glfw.KeyRelease
}
func (this *Listener) KeyPressed(key int) bool { //just once
	if glfw.Key(key) == glfw.KeyRelease {return false}
	_,found := this.mKeys[key]
	if !found {
		this.mKeys[key] = keyDown
		return true
	}
	if this.mKeys[key] == keyPressed {
		this.mKeys[key] = keyDown
		return true
	}
	return false
}

func (this *Listener) ButtonDown(key int) bool {
	return glfw.MouseButton(key) == glfw.KeyPress
}
func (this *Listener) ButtonUp(key int) bool {
	return glfw.MouseButton(key) == glfw.KeyRelease
}
func (this *Listener) ButtonPressed(key int) bool { //just once
	if glfw.MouseButton(key) == glfw.KeyRelease {return false}
	_,found := this.mButtons[key]
	if !found {
		this.mButtons[key] = keyDown
		return true
	}
	if this.mButtons[key] == keyPressed {
		this.mButtons[key] = keyDown
		return true
	}
	return false
}

func OnClose      () int       {return GListener.onClose()}
func OnMouseButton(button, state int) {GListener.onMouseButton(button, state)}
func OnMouseWheel (delta int)         {GListener.onMouseWheel(delta)}
func OnKey        (key, state int)    {GListener.onKey(key, state)}
func OnChar       (key, state int)    {GListener.onChar(key, state)}
//func onResize     (w, h int)          {GListener.onResize(w,h)}

func (this *Listener) onClose() int {
	fmt.Println("closed")
	this.quit = true
	return 0 // return 0 to keep window open.
}
func (this *Listener) onMouseButton(button, state int) {
	//fmt.Printf("mouse button: %d, %d\n", button, state)
	if state == glfw.KeyPress {
		this.mButtons[button] = keyPressed
	} else {
		this.mButtons[button] = keyUp
	}
}
func (this *Listener) onMouseWheel(delta int) {
	//fmt.Printf("mouse wheel: %d\n", delta)
}
func (this *Listener) onKey(key, state int) {
	//fmt.Printf("key: %d, %d\n", key, state)
	if state == glfw.KeyPress {
		this.mKeys[key] = keyPressed
	} else {
		this.mKeys[key] = keyUp
	}
}
func (this *Listener) onChar(key, state int) {
	//fmt.Printf("char: %d, %d\n", key, state)
	//this.onKey(key, state) //redirect
}
//func (this *Listener) onResize(w, h int) {
//	fmt.Println("resized:", w, h)
//}



