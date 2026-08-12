//go:build !js && !wasm

package main

import (
	"log"

	"github.com/gen2brain/malgo"
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

	device, err := malgo.InitDevice(ctx.Context, config, malgo.DeviceCallbacks{Data: debugConsole.Apu.MalgoAdapter})
	if err != nil {
		log.Fatal(err)
	}

	device.Start()
}
