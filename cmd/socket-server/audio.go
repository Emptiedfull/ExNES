package main

import (
	"fmt"
	"log"
	"math"

	"github.com/ebitengine/oto/v3"
)

var (
	audioCtx    *oto.Context
	audioPlayer *oto.Player
)

func setUpAudio() {
	ctx, readyChan, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 1,
		Format:       oto.FormatFloat32LE,
	})

	if err != nil {
		log.Fatal(err)
	}

	<-readyChan

	fmt.Println("pplaying audio")
	audioCtx = ctx // keep alive

	player := ctx.NewPlayer(debugConsole.Console.Apu.ExposedBuf)
	player.Play()

	audioPlayer = player

}

type sineReader struct {
	phase float64
}

func (s *sineReader) Read(p []byte) (int, error) {
	for i := 0; i+3 < len(p); i += 4 {
		sample := float32(math.Sin(s.phase * 2 * math.Pi * 440 / 44100))
		s.phase++
		bits := math.Float32bits(sample * 0.3)
		p[i] = byte(bits)
		p[i+1] = byte(bits >> 8)
		p[i+2] = byte(bits >> 16)
		p[i+3] = byte(bits >> 24)
	}
	return len(p), nil
}
