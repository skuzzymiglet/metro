package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/eiannone/keyboard"
	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
	"github.com/gabriel-vasile/mimetype"
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

func getReadCloser(name string) (io.ReadCloser, error) {
	var f io.ReadCloser
	var err error
	if _, err := pkger.Stat(name); err == nil {
		fmt.Println("embedded")
		f, err = pkger.Open(name)
	} else {
		fmt.Println("os.file")
		f, err = os.Open(name)
	}
	if err != nil {
		log.Fatal(err)
	}
	return f, err
}

func getStreamer(r io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	mime, err := mimetype.DetectReader(tee)
	if err != nil {
		log.Fatal(err)
	}
	var streamer beep.StreamSeekCloser
	var format beep.Format
	if mime.Is("audio/flac") {
		fmt.Println("flac")
		streamer, format, err = flac.Decode(ioutil.NopCloser(&buf))
	} else if mime.Is("audio/mp3") {
		fmt.Println("mp3")
		streamer, format, err = mp3.Decode(ioutil.NopCloser(&buf))
	} else if mime.Is("audio/wav") {
		streamer, format, err = wav.Decode(ioutil.NopCloser(&buf))
	} else if mime.Is("audio/ogg") {
		streamer, format, err = vorbis.Decode(ioutil.NopCloser(&buf))
	} else {
		log.Fatal("Sample is not flac, mp3, wav or ogg")
	}
	return streamer, format, err
}

func main() {
	pkger.Include("/samples/tabla_te_m.flac")

	fname := flag.String("f", "/samples/tabla_dhec.flac", "file")
	tempo := flag.Int("t", 120, "tempo")
	beats := flag.Int("b", 4, "beats")
	onSymbol := flag.String("o", "🔴", "Symbol for current beat")
	offSymbol := flag.String("O", "⭕", "Symbol for all other beats")
	flag.Parse()

	reader, err := getReadCloser(*fname)
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := getStreamer(reader)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(streamer.Len())

	buffer := beep.NewBuffer(format)
	fmt.Println(buffer.Len())
	toAppend := toTempo(streamer, format, *tempo)
	buffer.Append(toAppend)
	streamer.Close()
	fmt.Println(buffer.Len())
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
		fmt.Println("ree")
		speaker.Play(beep.Seq(loop))
		fmt.Println("ree")
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
