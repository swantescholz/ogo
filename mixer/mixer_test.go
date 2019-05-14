package mixer

import (
	"fmt"
	"ogo/mixer/music"
	"ogo/mixer/util"
	"ogo/mixer/sound"
	"ogo/mixer/timefunction"
	"math"
	"testing"
)

func testMusic() {
	mus := music.Load("music")
	//mus.FadeInPos(5.4, 8.8)
	mus.FadeIn(3.2)
	util.Wait(5)
	music.Pause()
	util.Wait(2)
	music.Resume()
	util.Wait(4)
	music.SetVolume(0.2)
	util.Wait(4)
	music.SetVolume(1.0)
	util.Wait(4)
	music.FadeOut(8.0)
	util.Wait(8)
}
func testSound(s *sound.Sound) {
	s.FadeIn(3.0, -1)
	util.Wait(5)
	s.Pause()
	util.Wait(2)
	s.Resume()
	util.Wait(4)
	s.SetVolume(0.2)
	util.Wait(4)
	s.SetVolume(1.0)
	util.Wait(4)
	s.FadeOut(8.0)
	util.Wait(6.5)
	sound.HaltAll()
}
func testSoundData() {
	sps := 22050
	dur := 34.0
	l := int(float64(sps)*dur)
	data := make([]float64, l)
	linv := 1.0/float64(l-1)
	for i := 0; i < l; i++ {
		f := float64(i)*linv*dur
		data[i] = 0.31*math.Sin(f*888.0)
	}
	sound.LoadData("data", data, sps)
	testSound(sound.Get("data"))
}
func testTimeFunction(fun string) {
	sps := 22050
	dur := 34.0
	tf := timefunction.Create(fun)
	s := sound.LoadData("tf", tf.ToFloats(sps, dur, 0.0), sps)
	s.Play(-1)
	util.Wait(dur)
	s.Halt()
}

func TestAll(t *testing.T) {
	p := func (s string) {t.Log(s);fmt.Println(s)}
	b := func() {p("Done.")}
	p("Testing All parts of mixer package")
	p("Init...")
	Init(true, false, 22050, 2, 4096)
	b();
	music.SetDefaultPath("res/music/", ".mp3")
	sound.SetDefaultPath("res/sounds/", ".wav")
	
	p("Testing Music...")
	testMusic()
	b(); p("Testing (owl) sound...")
	testSound(sound.Load("sound"))
	b(); p("Testing sine-wave sound...")
	testSoundData()
	b(); p("Testing TimeFunction based complex sine-wave sound...")
	testTimeFunction("0.2*(sin(888*t)+sin(2*t)*sin(555*t)*2)")
	b(); p("Quit...")
	Quit()
	b()
}
