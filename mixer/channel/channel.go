package channel

import (
	
)

/*
#cgo LDFLAGS: -lSDL -lSDL_mixer
#include <SDL/SDL.h>
#include <SDL/SDL_mixer.h>
*/
import "C"

var gMaxChannels = Channel(-1)

type Channel int
func Allocate(max int) {
	if gMaxChannels != -1 {
		panic("multiple allocation")
	}
	C.Mix_AllocateChannels(C.int(max))
	gMaxChannels = Channel(max)
}
func GetFree() Channel {
	var c Channel
	for c = 0; c < gMaxChannels; c++ {
		if !Playing(c) && !Paused(c) {return c}
	}
	panic("no free channel available")
}
func Playing(c Channel) bool {return C.Mix_Playing(C.int(c))==1}
func Paused (c Channel) bool {return C.Mix_Paused (C.int(c))==1}


