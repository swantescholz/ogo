package sound

import (
	"fmt"
	"log"
	"ogo/mixer/util"
	"unsafe"
	. "ogo/glmath/vec3"
	. "ogo/common"
)

//#cgo LDFLAGS: -lSDL -lSDL_mixer
//#include <SDL/SDL.h>
//#include <SDL/SDL_mixer.h>
import "C"

var GDefaultVolume Double = 0.7

type Sound struct {
	name string
	chunk *C.Mix_Chunk
	channel int
	defaultVolume, volumeFactor3D, volume Double
}
func New(name string, chunk *C.Mix_Chunk) *Sound {
	m := new(Sound)
	m.name = name
	m.chunk = chunk
	m.channel = -1
	m.defaultVolume, m.volumeFactor3D, m.volume = 1.0,1.0,1.0
	return m
}
func (t *Sound) Copy() *Sound {
	return New(t.name, t.chunk)
}
func (t *Sound) Destroy() {
	C.Mix_FreeChunk(t.chunk)
	t.chunk = nil
}

func (t *Sound) SetVolume(volume Double) {
	t.defaultVolume = volume
	t.updateVolume()
}
func (t *Sound) SetVolume3D(spos, lpos *Vec3) {
	t.volumeFactor3D = util.ComputeVolume3D(spos, lpos)
	t.updateVolume()
}
func (t *Sound) updateVolume() {
	t.volume = GDefaultVolume * t.defaultVolume * t.volumeFactor3D
	t.volume = t.volume.Clamp01()
	C.Mix_Volume(C.int(t.channel), C.int(t.volume*C.MIX_MAX_VOLUME))
}

func (t *Sound) Play(loops int) {
	t.channel = int(C.Mix_PlayChannelTimed(C.int(-1), t.chunk, C.int(loops), C.int(-1)))
	if t.channel == -1 {
		panic(fmt.Sprintf("Unable to play Sound file (%v): %v", t.name, util.GetMixError()))
	}
	t.SetVolume(GDefaultVolume)
}

func (t *Sound) FadeIn(duration Double, loops int) {
	t.channel = int(C.Mix_FadeInChannelTimed(C.int(-1), t.chunk, C.int(loops), C.int(duration*1000.0), C.int(-1)))
	if t.channel == -1 {
		panic(fmt.Sprintf("Unable to FadeIn Sound file (%v): %v", t.name, util.GetMixError()))
	}
	t.SetVolume(GDefaultVolume)
}
func (t *Sound) FadeOut(duration Double) {
	C.Mix_FadeOutChannel(C.int(t.channel), C.int(duration*1000.0))
}
func (t *Sound) Pause() {
	C.Mix_Pause(C.int(t.channel))
}
func (t *Sound) Resume() {
	C.Mix_Resume(C.int(t.channel))
}
func (t *Sound) Halt() {
	C.Mix_HaltChannel(C.int(t.channel))
}
func HaltAll() {C.Mix_HaltChannel(-1)}

var (
	GManager map[string]*C.Mix_Chunk
	GPathStart, GPathEnd string
)
func init() {
	GManager = make(map[string]*C.Mix_Chunk)
	GPathStart = "res/sounds/"
	GPathEnd = ".wav"
}
func SetDefaultPath(s1, s2 string) {
	GPathStart = s1
	GPathEnd = s2
}

func loadSound(name string) *C.Mix_Chunk {
	cfile, rb := C.CString(name), C.CString("rb")
	rwop := C.SDL_RWFromFile(cfile, rb)
	C.free(unsafe.Pointer(cfile))
	C.free(unsafe.Pointer(rb))
	chunk := C.Mix_LoadWAV_RW(rwop, 1)
	if chunk == nil {
		panic(fmt.Sprintf("Unable to load Sound file (%v): %v", name, util.GetMixError()))
	}
	return chunk
}

func loadSoundFromData(data []Double, sps int) *C.Mix_Chunk {
	bytes := util.CreateWavBytes(data, sps)
	rwop := C.SDL_RWFromConstMem(unsafe.Pointer(&bytes[0]), C.int(len(bytes)))
	chunk := C.Mix_LoadWAV_RW(rwop, 1)
	if chunk == nil {
		panic(fmt.Sprintf("Unable to load Sound data: %v",  util.GetMixError()))
	}
	return chunk
}
func completeName(name string) string {
	return GPathStart + name + GPathEnd
}
func Get(name string) *Sound {
	name = completeName(name)
	v := GManager[name]
	if v == nil {
		panic(fmt.Sprintf("Sound to get does not exist: %v", name))
	}
	return New(name, v)
}
func Load(name string) *Sound {
	name = completeName(name)
	v := GManager[name]
	if v != nil {
		log.Printf("Warning: Sound to load does already exist: %v\n", name)
		return New(name, v)
	}
	v = loadSound(name)
	GManager[name] = v
	return New(name, v)
}
func LoadData(name string, data []Double, sps int) *Sound {
	name = completeName(name)
	v := GManager[name]
	if v != nil {
		log.Printf("Warning: Sound data to load does already exist: %v\n", name)
		return New(name, v)
	}
	v = loadSoundFromData(data, sps)
	GManager[name] = v
	return New(name, v)
}
func Destroy(name string) {
	name = completeName(name)
	v := GManager[name]
	if v == nil {
		log.Printf("Warning: Sound to delete does not exist: %v\n", name)
		return
	}
	New(name, v).Destroy()
	delete(GManager, name)
}
func ClearAllResources() {
	for name,val := range GManager {
		New(name, val).Destroy()
	}
	GManager = nil
	GManager = make(map[string]*C.Mix_Chunk)
}


//*/
