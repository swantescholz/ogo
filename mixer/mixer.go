package mixer

import (
	"fmt"
	"ogo/mixer/util"
	"ogo/mixer/channel"
)

/*
#cgo LDFLAGS: -lSDL -lSDL_mixer
#include <SDL/SDL.h>
#include <SDL/SDL_mixer.h>
*/
import "C"

var GSDLWasInitHere = false
func Init(initSDL, stereo bool, rate, nchannels, nbuffers int) {
	if initSDL {
		GSDLWasInitHere = true
		if C.SDL_Init(C.SDL_INIT_AUDIO) != 0 {
			panic(fmt.Sprintf("Unable to initialize SDL: %v\n", util.GetSdlError()))

		}
	}
	
	//initFlags := C.MIX_INIT_FLAC | C.MIX_INIT_MP3 | C.MIX_INIT_OGG
	//C.Mix_Init(initFlags)
	
	audio_format := C.AUDIO_S16SYS
	nstereo := 1
	if stereo {nstereo = 2}
	if C.Mix_OpenAudio(C.int(rate), C.Uint16(audio_format),
		C.int(nstereo), C.int(nbuffers)) != 0 {
		panic(fmt.Sprintf("Unable to initialize audio: %v\n", util.GetMixError()))
	}
	channel.Allocate(nchannels)
}

func Quit() {
	C.Mix_CloseAudio()	
	//for C.Mix_Init(0) != 0 {C.Mix_Quit()} // force a quit
	if GSDLWasInitHere {C.SDL_Quit()}
}


