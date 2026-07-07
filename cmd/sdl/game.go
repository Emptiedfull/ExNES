package main

import (
	"exnes/Core"
	"fmt"
	"os"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type game struct {
	core        *Core.Console
	audioDevice sdl.AudioDeviceID

	pauseChannel chan bool
}

func (console *game) LoadRom(filepath string, mb *menuBar) error {
	file, err := os.Open(filepath)

	if err != nil {
		return fmt.Errorf("uhm: %v", err)
	}

	err = console.core.InitRom(file)
	console.core.Cpu.Reset()
	fmt.Println("error:", err)
	fmt.Println("starting the game")

	go beginSampleLoop(console.audioDevice, console.core, console.pauseChannel)
	sdl.PauseAudioDevice(console.audioDevice, false)

	mb.setFlag(gameRunning, true)
	mb.setFlag(gamePlaying, true)

	return nil
}

func initConsole(screenChannel chan Core.ScreenInfo) *Core.Console {
	console := Core.InitializeConsole()

	console.Apu = Core.NewApu(44100, console)
	console.ScreenChannel = screenChannel

	return console
}

func frameMonitor(console *Core.Console) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	last := 0

	for range ticker.C {
		framesRun := console.Ppu.Frame - last

		fmt.Println("fps:", framesRun)

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
		sdl.QueueAudio(dev, float32ToBytes(sampleBuf))
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
