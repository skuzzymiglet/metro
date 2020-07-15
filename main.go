package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	"errors"

	"github.com/eiannone/keyboard"
	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
	"github.com/markbates/pkger"

	"log"
	"os"
	"strings"
	"time"
)

func beatString(beat, beats int, on, off string) string {
	s := make([]string, beats)
	for i := 0; i < beats; i++ {
		if i == beat {
			s[i] = on
		} else {
			s[i] = off
		}
	}
	return strings.Join(s, "")
}

func toTempo(s beep.StreamSeeker, f beep.Format, tempo int) *beep.Resampler {
	ratio := float64(s.Len()) / float64(f.SampleRate.N(time.Duration(60000/tempo)*time.Millisecond))
	return beep.ResampleRatio(1, ratio, s)
	// return beep.ResampleRatio(2, , s)
}

func getFile(name string) (io.ReadSeeker, error) {
	var f io.ReadSeeker
	var err error
	if _, err := pkger.Stat(name); err == nil {
		fmt.Printf("Playing embedded sample %s\n", name)
		f, err = pkger.Open(name)
	} else {
		f, err = os.Open(name)
	}
	if err != nil {
		log.Fatal(err)
	}
	return f, err
}

func getStreamer(r io.ReadSeeker) (beep.StreamSeekCloser, beep.Format, error) {
	var streamer beep.StreamSeekCloser
	var format beep.Format
	nopcloser := ioutil.NopCloser(r)
	r.Seek(0, 0)
	streamer, format, err := flac.Decode(nopcloser)
	fmt.Println("type flac")
	if err == nil {
		return streamer, format, err
	}
	r.Seek(0, 0)
	streamer, format, err = mp3.Decode(nopcloser)
	fmt.Println("type mp3")
	if err == nil {
		return streamer, format, err
	}
	r.Seek(0, 0)
	streamer, format, err = wav.Decode(nopcloser)
	fmt.Println("type wav")
	if err == nil {
		return streamer, format, err
	}
	streamer, format, err = vorbis.Decode(nopcloser)
	fmt.Println("type vorbis")
	if err == nil {
		return streamer, format, err
	}
	return streamer, format, errors.New("Invalid format - needs to be flac/mp3/wav/vorbis")
}

func main() {
	pkger.Include("/samples/tabla_te_m.flac")

	fname := flag.String("f", "/samples/tabla_te_m.flac", "file")
	tempo := flag.Int("t", 120, "tempo")
	beats := flag.Int("b", 4, "beats")
	onSymbol := flag.String("o", "🔴", "Symbol for current beat")
	offSymbol := flag.String("O", "⭕", "Symbol for all other beats")
	flag.Parse()

	reader, err := getFile(*fname)
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := getStreamer(reader)
	if err != nil {
		log.Fatal(err)
	}

	buffer := beep.NewBuffer(format)
	toAppend := toTempo(streamer, format, *tempo)
	buffer.Append(toAppend)
	streamer.Close()
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
	go func() {
		speaker.Play(beep.Seq(loop))
	}()
	for {
		c := time.Tick(time.Duration(60000 / *tempo) * time.Millisecond)
		currentBeat := 0
		fmt.Printf("\r%s", beatString(currentBeat, *beats, *onSymbol, *offSymbol))
		currentBeat++
		currentBeat %= *beats
		for range c {
			go func() {
				fmt.Printf("\r%s", beatString(currentBeat, *beats, *onSymbol, *offSymbol))
				currentBeat++
				currentBeat %= *beats
			}()
		}
	}
}
