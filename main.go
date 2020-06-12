package main

import (
	// "flag"
	"fmt"
	// "github.com/eiannone/keyboard"
	// "github.com/faiface/beep"
	// "github.com/faiface/beep/flac"
	// "github.com/faiface/beep/speaker"
	"github.com/go-mix/mix"
	"github.com/go-mix/mix/bind"
	"github.com/go-mix/mix/bind/spec"
	// "log"
	"os"
	// "strings"
	"math/rand"
	"time"
)

var (
	loader, out string
	profileMode string
	sampleHz    = float64(48000)
	specs       = spec.AudioSpec{
		Freq:     sampleHz,
		Format:   spec.AudioF32,
		Channels: 2,
	}
	bpm     = 120
	step    = time.Minute / time.Duration(bpm*4)
	loops   = 8
	prefix  = "samples/"
	kick1   = "samples/drum_heavy_kick.flac"
	kick2   = "samples/drum_tom_mid_hard.flac"
	pattern = []string{
		kick1,
		kick2,
		kick1,
		kick2,
	}
)

// func beatString(beat, beats int) string {
// 	s := make([]string, beats)
// 	for i := 0; i < beats; i++ {
// 		if i == beat {
// 			s[i] = "🔴"
// 		} else {
// 			s[i] = "⭕"
// 		}
// 	}
// 	return strings.Join(s, "")
// }

func main() {
	out := os.Stdout
	bind.UseOutputString(out)
	bind.UseLoaderString(loader)
	defer mix.Teardown()
	mix.Configure(specs)
	mix.SetSoundsPath(prefix)

	// setup the music
	t := 1 * time.Second // buffer before music
	for n := 0; n < loops; n++ {
		for s := 0; s < len(pattern); s++ {
			mix.SetFire(
				pattern[s], t+time.Duration(s)*step, 0, 1.0, rand.Float64()*2-1)
		}
		t += time.Duration(len(pattern)) * step
	}
	t += 5 * time.Second // buffer after music

	//
	if bind.IsDirectOutput() {
		out := os.Stdout
		mix.Debug(true)
		mix.OutputStart(t, out)
		for p := time.Duration(0); p <= t; p += t / 4 {
			mix.OutputContinueTo(p)
		}
		mix.OutputClose()
	} else {
		mix.Debug(true)
		mix.StartAt(time.Now().Add(1 * time.Second))
		fmt.Printf("Mix: 808 Example - pid:%v playback:%v spec:%v\n", os.Getpid(), out, specs)
		for mix.FireCount() > 0 {
			time.Sleep(1 * time.Second)
		}
	}
	// keysEvents, err := keyboard.GetKeys(10)
	// if err != nil {
	// 	panic(err)
	// }
	// defer func() {
	// 	_ = keyboard.Close()
	// }()
	// go func() {
	// 	for {
	// 		keyEvent := <-keysEvents
	// 		if keyEvent.Key == keyboard.KeyEsc {
	// 			fmt.Println("Exiting, ESC recieved")
	// 			os.Exit(0)
	// 		}
	// 	}
	// }()
	// fname := flag.String("f", "samples/tabla_te2.flac", "file")
	// tempo := flag.Float64("t", 120, "tempo")
	// beats := flag.Int("b", 4, "beats")
	// flag.Parse()
	// f, err := os.Open(*fname)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// streamer, format, err := flac.Decode(f)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	// buffer := beep.NewBuffer(format)
	// buffer.Append(streamer)
	// streamer.Close()
	// for {
	// 	// char, key, err := keyboard.GetKey()
	// 	// if err != nil {
	// 	// 	panic(err)
	// 	// }
	// 	// fmt.Println(key, char)
	// 	// if key == keyboard.KeyEsc {
	// 	// 	break
	// 	// } else if key == keyboard.KeySpace {
	// 	c := time.Tick(time.Duration(60.0 / *tempo * 1000000000.0) * time.Nanosecond)
	// 	currentBeat := 0
	// 	s := time.Now()
	// 	for tick := range c {
	// 		beat := buffer.Streamer(0, buffer.Len())
	// 		go func() {
	// 			speaker.Play(beat)
	// 			fmt.Printf("\r%s at %s", beatString(currentBeat, *beats), tick)
	// 			currentBeat++
	// 			currentBeat %= *beats
	// 			fmt.Printf("delay: %s", time.Now().Sub(s))
	// 		}()
	// 	}
	// 	// }
	// }
}
