package listener


import (
	"github.com/banthar/Go-SDL/sdl"
)

const (
	keyUp = iota
	keyPressed = iota
	keyDown = iota
)

var GListener *Listener = nil

type Listener struct {
	quit bool
	keys map[uint32]uint32
	//buttons map[uint32]uint32
}
func New() *Listener {
	il := &Listener{false, make(map[uint32]uint32)}
	return il
}

func (this *Listener) Update() {
	for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
		switch e.(type) {
			case *sdl.QuitEvent:
				this.quit = true
				break;
			case *sdl.KeyboardEvent:
				ke := e.(*sdl.KeyboardEvent);
				if ke.Type == sdl.KEYDOWN {
					this.keys[ke.Keysym.Sym] = keyPressed
				} else if ke.Type == sdl.KEYUP {
					this.keys[ke.Keysym.Sym] = keyUp
				}
		}
	}
}

func (this *Listener) MbQuit() bool {
	return this.quit
}

func (this *Listener) KeyDown(key uint32) bool {
	value,found := this.keys[key]
	if !found {return false}
	if value == keyDown || value == keyPressed {
		return true
	}
	return false
}
func (this *Listener) KeyUp(key uint32) bool {
	return !this.KeyDown(key)
}
func (this *Listener) KeyPressed(key uint32) bool { //just once
	_,found := this.keys[key]
	if !found {
		return false
	}
	if this.keys[key] == keyPressed {
		this.keys[key] = keyDown
		return true
	}
	return false
}

/*
func (this *Listener) ButtonDown(key uint32) bool {
	return glfw.MouseButton(key) == glfw.KeyPress
}
func (this *Listener) ButtonUp(key uint32) bool {
	return glfw.MouseButton(key) == glfw.KeyRelease
}
func (this *Listener) ButtonPressed(key uint32) bool { //just once
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
}//*/



