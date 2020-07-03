package main

import (
	"flag"
	"fmt"
	"io"

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

func main() {
	pkger.Include("/samples/tabla_te_m.flac")

	fname := flag.String("f", "samples/tabla_dhec.flac", "file")
	tempo := flag.Int("t", 120, "tempo")
	beats := flag.Int("b", 4, "beats")
	onSymbol := flag.String("o", "🔴", "Symbol for current beat")
	offSymbol := flag.String("O", "⭕", "Symbol for all other beats")
	flag.Parse()
	var f io.ReadCloser
	if _, err := pkger.Stat(*fname); err == nil {
		fmt.Println("embedded")
		f, err = pkger.Open(*fname)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("os.file")
		f, err = os.Open(*fname)
		if err != nil {
			log.Fatal(err)
		}
	}
	var streamer beep.StreamSeekCloser
	var format beep.Format
	b := make([]byte, 64)
	_, err := f.Read(b)
	// fmt.Println(f.Stat().Size())
	fmt.Printf("%T %#v\n", f, f)
	mime, err := mimetype.DetectReader(f)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(mime.String())
	fmt.Println(f)
	if mime.Is("audio/flac") {
		fmt.Println("flac")
		streamer, format, err = flac.Decode(f)
	} else if mime.Is("audio/mp3") {
		fmt.Println("mp3")
		streamer, format, err = mp3.Decode(f)
	} else if mime.Is("audio/wav") {
		streamer, format, err = wav.Decode(f)
	} else if mime.Is("audio/ogg") {
		streamer, format, err = vorbis.Decode(f)
	} else {
		log.Fatal("Sample is not flac, mp3, wav or ogg")
	}
	fmt.Println(io.LimitReader(f, 100))
	if err != nil {
		panic(err)
	}
	// streamer, format, err := mp3.Decode(f)
	// if err != nil {
	// 	fmt.Println("not mp3")
	// 	streamer, format, err = flac.Decode(f)
	// 	log.Println(err)
	// 	if err != nil {
	// 		fmt.Println("not flac")
	// 		streamer, format, err = wav.Decode(f)
	// 		log.Println(err)
	// 		if err != nil {
	// 			fmt.Println("not wav")
	// 			streamer, format, err = vorbis.Decode(f)
	// 			log.Println(err)
	// 			if err != nil {
	// 				fmt.Println("not vorbis")
	// 				log.Println("Sample is none of supported formats: mp3, flac, wav, vorbis")
	// 				log.Fatal(err)
	// 			}
	// 		}
	// 	}
	// }

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
