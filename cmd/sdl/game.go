package main

import (
	"exnes/Core"
	"fmt"
	"log"
	"os"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type game struct {
	core        *Core.Console
	audioDevice sdl.AudioDeviceID
	romPath     string

	pauseChannel  chan bool
	screenChannel chan Core.ScreenInfo

	fps int
}

func (console *game) LoadRom(filepath string, mb *menuBar) error {
	file, err := os.Open(filepath)

	if err != nil {
		return fmt.Errorf("uhm: %v", err)
	}

	err = console.core.InitRom(file)
	console.core.Cpu.Reset()
	fmt.Println("error:", err)

	go beginSampleLoop(console.audioDevice, console.core, console.pauseChannel)
	go console.frameMonitor(console.core)
	sdl.PauseAudioDevice(console.audioDevice, false)

	console.romPath = filepath

	mb.setFlag(gameRunning, true)
	mb.setFlag(gamePlaying, true)

	return nil
}

func (g *game) initConsole(screenChannel chan Core.ScreenInfo) {
	g.core = Core.InitializeConsole()

	g.core.Apu = Core.NewApu(44100, g.core)
	g.core.ScreenChannel = screenChannel
	g.screenChannel = screenChannel

}

func (g *game) reloadROM() {

	g.core.PowerCycle()

	file, err := os.Open(g.romPath)

	if err != nil {
		log.Fatal("AHHHH:", err)
	}

	g.core.InitRom(file)
	g.core.Cpu.Reset()

}

func (g *game) frameMonitor(console *Core.Console) {

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	last := 0

	for range ticker.C {
		if console == nil {
			continue
		}
		framesRun := console.Ppu.Frame - last

		g.fps = framesRun

		last = console.Ppu.Frame
	}
}

const bufferMillis = 46

func backpressureThreshold(freq int32) uint32 {
	bytesPerSecond := uint32(freq) * 4
	return bytesPerSecond * bufferMillis / 1000
}

func beginSampleLoop(dev sdl.AudioDeviceID, console *Core.Console, pauseChannel chan bool) {
	const samples = 512

	sampleBuf := make([]float32, samples)
	threshold := backpressureThreshold(44100)

	for {

		select {
		case paused := <-pauseChannel:
			if paused {
				for state := range pauseChannel {
					if !state {
						break
					}
				}
			}
		default:
		}

		for i := range samples {
			for !console.Apu.HasSample() {
				console.TickNoAudio()

			}
			sample := console.Apu.PopSample()
			sampleBuf[i] = sample

		}
		for sdl.GetQueuedAudioSize(dev) > threshold {
			sdl.Delay(1)
		}
		console.RunDisplayUpdates()
		adjustVolume(sampleBuf, 0.1)
		sdl.QueueAudio(dev, float32ToBytes(sampleBuf))
	}

}

func adjustVolume(samples []float32, volume float32) {
	if volume == 1 {
		return
	}

	for i := range samples {
		samples[i] = samples[i] * volume
	}

}

func float32ToBytes(samples []float32) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&samples[0])), len(samples)*4)
}

var controlMap = map[sdl.Keycode]int{
	sdl.K_z:      Core.ButtonA,
	sdl.K_x:      Core.ButtonB,
	sdl.K_UP:     Core.ButtonUp,
	sdl.K_DOWN:   Core.ButtonDown,
	sdl.K_LEFT:   Core.ButtonLeft,
	sdl.K_RIGHT:  Core.ButtonRight,
	sdl.K_LSHIFT: Core.ButtonSelect,
	sdl.K_RETURN: Core.ButtonStart,
}

func handleInputs(console *Core.Console, e *sdl.KeyboardEvent) {

	pressed := e.State == sdl.PRESSED

	console.Player1.UpdateBtnBool(controlMap[e.Keysym.Sym], pressed)
}

func (console *game) changeSpeed(multipler float64) {

	if multipler > 1 {
		desiredfps := 60 * multipler
		internal := desiredfps / 40
		multipler = internal
	}

	targetFrequency := 44100 / multipler
	console.core.Apu = Core.NewApu(targetFrequency, console.core)

}
