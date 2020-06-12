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

func main() {
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
			if keyEvent.Key == keyboard.KeyEsc {
				fmt.Println("Exiting, ESC recieved")
				os.Exit(0)
			}
		}
	}()
	fname := flag.String("f", "samples/tabla_te2.flac", "file")
	tempo := flag.Float64("t", 120, "tempo")
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
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()
	for {
		// char, key, err := keyboard.GetKey()
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Println(key, char)
		// if key == keyboard.KeyEsc {
		// 	break
		// } else if key == keyboard.KeySpace {
		c := time.Tick(time.Duration(60.0 / *tempo * 1000000000.0) * time.Nanosecond)
		currentBeat := 0
		s := time.Now()
		for tick := range c {
			beat := buffer.Streamer(0, buffer.Len())
			go func() {
				speaker.Play(beat)
				fmt.Printf("\r%s at %s", beatString(currentBeat, *beats), tick)
				currentBeat++
				currentBeat %= *beats
				fmt.Printf("delay: %s", time.Now().Sub(s))
			}()
		}
		// }
	}
}
