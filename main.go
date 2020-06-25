package main

import (
	"flag"
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/speaker"
	"log"
	"os"
	"strings"
	"time"
)

func beatString(beat, beats int) string {
	s := make([]string, beats)
	for i := 0; i < beats; i++ {
		if i == beat {
			s[i] = "🔴"
		} else {
			s[i] = "⭕"
		}
	}
	return strings.Join(s, "")
}

func toTempo(s beep.StreamSeeker, f beep.Format, tempo int) *beep.Resampler {
	ratio := float64(s.Len()) / float64(f.SampleRate.N(time.Duration(60000/tempo)*time.Millisecond))
	return beep.ResampleRatio(1, ratio, s)
	// return beep.ResampleRatio(2, , s)
}

func main() {
	fname := flag.String("f", "samples/tabla_te2.flac", "file")
	tempo := flag.Int("t", 120, "tempo")
	beats := flag.Int("b", 4, "beats")
	flag.Parse()
	f, err := os.Open(*fname)
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := flac.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	buffer := beep.NewBuffer(format)
	toAppend := toTempo(streamer, format, *tempo)
	buffer.Append(toAppend)
	loop := beep.Loop(-1, buffer.Streamer(0, buffer.Len()))
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()
	go func() {
		for {
			keyEvent := <-keysEvents
			if keyEvent.Rune == 'q' {
				fmt.Print("\r")
				os.Exit(0)
			}
		}
	}()
	go speaker.Play(beep.Seq(loop))
	for {
		c := time.Tick(time.Duration(60000 / *tempo) * time.Millisecond)
		currentBeat := 0
		fmt.Printf("\r%s", beatString(currentBeat, *beats))
		currentBeat++
		currentBeat %= *beats
		for range c {
			go func() {
				fmt.Printf("\r%s", beatString(currentBeat, *beats))
				currentBeat++
				currentBeat %= *beats
			}()
		}
	}
}
