package music

import (
	"fmt"
	"log"
	"ogo/mixer/util"
	"unsafe"
)

/*
#cgo LDFLAGS: -lSDL -lSDL_mixer
#include <SDL/SDL.h>
#include <SDL/SDL_mixer.h>
*/
import "C"

type Music struct {
	name string
	mus *C.Mix_Music
}
func New(name string, music *C.Mix_Music) *Music {
	m := new(Music)
	m.name = name
	m.mus = music
	return m
}
func (t *Music) Destroy() {
	C.Mix_FreeMusic(t.mus)
	t.mus = nil
}

func (t *Music) Play() {
	t.PlayLoops(0)
}
func (t *Music) PlayLoops(loops int) {
	if C.Mix_PlayMusic(t.mus, C.int(loops)) == -1 {
		panic(fmt.Sprintf("Unable to play Music file (%v): %v", t.name, util.GetMixError()))
	}
}
func (t *Music) FadeIn(duration float64) {
	t.FadeInLoops(duration, 0)
}
func (t *Music) FadeInLoops(duration float64, loops int) {
	if C.Mix_FadeInMusic(t.mus, C.int(loops), C.int(duration*1000.0)) == -1 {
		panic(fmt.Sprintf("Unable to FadeIn Music file (%v): %v", t.name, util.GetMixError()))
	}
}
func (t *Music) FadeInPos(duration,pos float64) {
	t.FadeInLoopsPos(duration, 0, pos)
}
func (t *Music) FadeInLoopsPos(duration float64, loops int, pos float64) {
	if C.Mix_FadeInMusicPos(t.mus, C.int(loops), C.int(duration*1000.0), C.double(pos)) == -1 {
		panic(fmt.Sprintf("Unable to FadeIn Music file (%v): %v", t.name, util.GetMixError()))
	}
}
var (
	GManager map[string]*Music
	GPathStart, GPathEnd string
)
func init() {
	GManager = make(map[string]*Music)
	GPathStart = "res/music/"
	GPathEnd = ".ogg"
}
func SetDefaultPath(s1, s2 string) {
	GPathStart = s1
	GPathEnd = s2
}

func loadMusic(name string) *Music {
	cname := C.CString(name)
	music := C.Mix_LoadMUS(cname)
	C.free(unsafe.Pointer(cname))
	if music == nil {
		panic(fmt.Sprintf("Unable to load Music file (%v): %v", name, util.GetMixError()))
	}
	return New(name, music)
}
func completeName(name string) string {
	return GPathStart + name + GPathEnd
}
func Get(name string) *Music {
	name = completeName(name)
	v := GManager[name]
	if v == nil {
		panic(fmt.Sprintf("Music to get does not exist: %v", name))
	}
	return v
}
func Load(name string) *Music {
	name = completeName(name)
	v := GManager[name]
	if v != nil {
		log.Printf("Warning: Music to load does already exist: %v\n", name)
		return v
	}
	v = loadMusic(name)
	GManager[name] = v
	return v
}
func Destroy(name string) {
	name = completeName(name)
	v := GManager[name]
	if v == nil {
		log.Printf("Warning: Music to delete does not exist: %v\n", name)
		return
	}
	delete(GManager, name)
}
func ClearAllResources() {
	for _,val := range GManager {
		val.Destroy()
	}
	GManager = nil
	GManager = make(map[string]*Music)
}

func GetVolume() float64 {
	return float64(C.Mix_VolumeMusic(-1))/float64(C.MIX_MAX_VOLUME)
}
func SetVolume(volume float64) {
	C.Mix_VolumeMusic(C.int(volume*float64(C.MIX_MAX_VOLUME)))
}

func SetPosition(pos float64) {
	t := C.Mix_GetMusicType(nil)
	if t == C.MUS_MP3 {Rewind()}
	C.Mix_SetMusicPosition(C.double(pos))
}

func FadeOut(duration float64) {
	C.Mix_FadeOutMusic(C.int(duration*1000.0))
}

func Playing() bool {return C.Mix_PlayingMusic()==1}
func Paused () bool {return C.Mix_PausedMusic ()==1}
func Fading () int {
	f := C.Mix_FadingMusic()
	if f == C.MIX_FADING_IN  {return -1}
	if f == C.MIX_FADING_OUT {return  1}
	return 0
}
func Rewind() {C.Mix_RewindMusic()}
func Pause () {C.Mix_PauseMusic()}
func Resume() {C.Mix_ResumeMusic()}
func Halt  () {C.Mix_HaltMusic()}






