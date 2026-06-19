package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/gen2brain/malgo"
)

var (
	audioCtx    *oto.Context
	audioPlayer *oto.Player
)

func setUpAudioDriver() {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		log.Fatal(err)
	}

	config := malgo.DefaultDeviceConfig(malgo.Playback)

	config.Playback.Format = malgo.FormatF32
	config.Playback.Channels = 1
	config.SampleRate = 44100
	config.PeriodSizeInMilliseconds = 8
	config.Alsa.NoMMap = 1

	onSamples := func(output, input []byte, frameCount uint32) {
		if frameCount == 0 || len(output) == 0 {
			fmt.Println(len(output), len(input))
			return
		}

		samplesNeeded := int(frameCount)

		for i := range samplesNeeded {
			for !debugConsole.Console.Apu.HasSample() {
				debugConsole.Console.TickNoAudio()
			}

			sample := debugConsole.Console.Apu.PopSample()

			bits := math.Float32bits(sample)
			output[i*4] = byte(bits)
			output[i*4+1] = byte(bits >> 8)
			output[i*4+2] = byte(bits >> 16)
			output[i*4+3] = byte(bits >> 24)
		}

	}

	device, err := malgo.InitDevice(ctx.Context, config, malgo.DeviceCallbacks{Data: onSamples})
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		i := 1
		for range ticker.C {
			fmt.Println(debugConsole.Console.Ppu.Frame / i)
			i++
		}
	}()

	device.Start()
}

func setUpAudio() {

	ctx, readyChan, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 1,
		Format:       oto.FormatFloat32LE,
		BufferSize:   time.Millisecond * 20,
	})

	if err != nil {
		log.Fatal(err)
	}

	<-readyChan

	fmt.Println("pplaying audio")
	audioCtx = ctx // keep alive

	player := ctx.NewPlayer(debugConsole.Console.Apu)
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
