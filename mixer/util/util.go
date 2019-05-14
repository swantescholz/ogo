package util

import (
	"fmt"
	"time"
	"io/ioutil"
	. "ogo/glmath/vec3"
	. "ogo/common"
)

/*
#cgo LDFLAGS: -lSDL -lSDL_mixer
#include <SDL/SDL.h>
#include <SDL/SDL_mixer.h>
*/
import "C"

//3D math
func ComputePanning(lpos, ldir, lup, spos *Vec3) Double {
	side := lup.Cross(ldir).Unit()
	sl := spos.Sub(lpos)
	x := sl.Dot(side)
	z := sl.Dot(ldir)
	angle := x.Atan2(z)
	pan := angle.Sin()
	return pan
}
func ComputeVolume3D(lpos, spos *Vec3) Double {
	d2 := lpos.Distance2(spos)
	return 1.0/d2
}

func GetMixError() string {
	return C.GoString(C.Mix_GetError())
}
func GetSdlError() string {
	return C.GoString(C.SDL_GetError())
}

func Wait(sec Double) {
	s := fmt.Sprintf("%vs", sec)
	dur,_ := time.ParseDuration(s)
	time.Sleep(dur)
}

//little endian!
func AppendInt16(bytes []byte, i int16) []byte {
	a := byte((i>> 0) % 256)
	b := byte((i>> 8) % 256)
	return append(bytes, a,b)
}
func AppendInt32(bytes []byte, i int32) []byte {
	a := byte((i>> 0) % 256)
	b := byte((i>> 8) % 256)
	c := byte((i>>16) % 256)
	d := byte((i>>24) % 256)
	return append(bytes, a,b,c,d)
}

func CreateWavBytes(data []Double, sps int) []byte {
	size := int32(44+len(data)*2)
	bytes := make([]byte, 0, size)
	bytes = append(bytes, "RIFF"...)
	bytes = AppendInt32(bytes, size-8)
	bytes = append(bytes, "WAVE"...)
	bytes = append(bytes, "fmt "...)
	bytes = AppendInt32(bytes, 16)
	bytes = AppendInt16(bytes, 1)
	bytes = AppendInt16(bytes, 1)
	bytes = AppendInt32(bytes, int32(sps)) //samples per second
	bytes = AppendInt32(bytes, int32(sps)*2) //bytes per second
	bytes = AppendInt16(bytes, 2) //bytes per sample
	bytes = AppendInt16(bytes, 16) //bits per sample
	bytes = append(bytes, "data"...)
	bytes = AppendInt32(bytes, int32(len(data)))
	const max = (65500.0*0.5)
	for i := range data {
		f := data[i]
		var x int16
		x = int16(f*max)
		bytes = AppendInt16(bytes, x)
	}
	return bytes
}
func WriteWav(filename string, data []Double, sps int) {
	bytes := CreateWavBytes(data, sps)
	if err := ioutil.WriteFile(filename, bytes, 0644); err != nil {
		panic(err)
	}
}


