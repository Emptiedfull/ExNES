package main

import (
	"exnes/Core"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type game struct {
	core        *Core.Console
	audioDevice sdl.AudioDeviceID
	romPath     string

	pauseChannel  chan bool
	screenChannel chan []uint32

	fps    int
	volume float32

	TurboState map[Core.BUTTON]chan bool
	TurboMux   sync.Mutex

	gameLoaded bool
}

func (console *game) LoadRom(filepath string, mb *menuBar) error {

	file, err := os.Open(filepath)

	if err != nil {
		return err
	}
	defer file.Close()

	err = console.core.InitRom(file)
	if err != nil {
		return err
	}
	console.core.Cpu.Reset()
	console.romPath = filepath

	if !console.gameLoaded {
		go console.beginSampleLoop()
		go console.frameMonitor(console.core)
		sdl.PauseAudioDevice(console.audioDevice, false)
	}

	console.core.GetHash()
	mb.state.NewSaves[console.core.GetHash()] = make([]romSave, 10)

	mb.setFlag(gameRunning, true)
	mb.setFlag(gamePlaying, true)

	console.gameLoaded = true

	return nil
}

func (g *game) initConsole(screenChannel chan []uint32) {
	g.core = Core.InitializeConsole()

	g.core.LoadRam = LoadRam

	g.core.Apu = Core.NewApu(44100, g.core)
	g.core.ScreenChannel = screenChannel
	g.screenChannel = screenChannel

}

func (g *game) SaveRam() {
	if g.core.GetMapper() == nil {
		return
	}

	data := g.core.GetRam()
	if data != nil {

		saveName := string(g.core.GetHash()) + ".sav"
		saveWithFail(saveName, data)

	} else {
		fmt.Println("ignoring", g.core.GetName())
	}
}

func saveWithFail(name string, data []uint8) {
	saveDir := filepath.Join(targetDir(), "saves")
	_, err := os.Stat(saveDir)
	if err != nil {
		os.MkdirAll(saveDir, 0755)
	}

	err = os.WriteFile(filepath.Join(saveDir, name), data, 0644)
	if err != nil {
		pushError("Save", err, false)

	}

}

func LoadRam(hash string, ram []uint8) {

	path := filepath.Join(targetDir(), "./saves", hash+".sav")

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	copy(ram, data)

}

func (g *game) reloadROM() {

	fmt.Println("saving this ")
	g.SaveRam()

	g.core.PowerCycle()

	file, err := os.Open(g.romPath)

	if err != nil {

		pushError("ROM", err, false)
		return

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

func (game *game) beginSampleLoop() {
	const samples = 512

	sampleBuf := make([]float32, samples)
	threshold := backpressureThreshold(44100)

	for {

		select {
		case paused := <-game.pauseChannel:
			if paused {
				for state := range game.pauseChannel {
					if !state {
						break
					}
				}
			}
		default:
		}

		for i := range samples {
			for !game.core.Apu.HasSample() {
				game.core.Step()

			}
			sample := game.core.Apu.PopSample()
			sampleBuf[i] = sample

		}
		for sdl.GetQueuedAudioSize(game.audioDevice) > threshold {
			sdl.Delay(1)
		}
		game.core.RunDisplayUpdates()
		adjustVolume(sampleBuf, game.volume)
		sdl.QueueAudio(game.audioDevice, float32ToBytes(sampleBuf))
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

func handleInputs(console *game, e *sdl.KeyboardEvent, inp Inputs, turboInp Inputs) {

	if action, ok := inp.KeyToAction[e.Keysym.Scancode]; ok {
		pressed := e.State == sdl.PRESSED

		console.core.Player1.UpdateBtnBool(action, pressed)
		return
	}

	if action, ok := turboInp.KeyToAction[e.Keysym.Scancode]; ok {

		pressed := e.State == sdl.PRESSED

		if e.Repeat != 0 {

			return
		}

		if pressed {
			go console.beginTurbo(action)
		} else {
			go console.stopTurbo(action)
		}
	}

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

func (console *game) changeVolume(volumeStr string) {
	volumeStr = volumeStr[0 : len(volumeStr)-1]

	volumeInt, err := strconv.Atoi(volumeStr)
	if err != nil {

		pushError("Volume", err, false)
	}

	volume := float32(volumeInt) / 100.0

	if volume > 1 {
		volume = 1
	}

	console.volume = volume
}

func renderFrame(texture *sdl.Texture, renderer *sdl.Renderer, buffer []uint32, gameRect *sdl.Rect) {
	texture.Update(nil, unsafe.Pointer(&buffer[0]), 256*4)
	renderer.Copy(texture, nil, gameRect)

}

type Inputs struct {
	KeyToAction map[sdl.Scancode]Core.BUTTON
	ActionToKey map[Core.BUTTON]sdl.Scancode
}

func (base *Inputs) AssignKey(key sdl.Scancode, action Core.BUTTON) {

	for k, a := range base.KeyToAction {
		if a == action {
			delete(base.KeyToAction, k)
		}
	}

	for a, k := range base.ActionToKey {
		if k == key {
			delete(base.ActionToKey, a)
		}
	}

	base.KeyToAction[key] = action
	base.ActionToKey[action] = key

}

func (base *Inputs) DumbReadable() {
	fmt.Println("Action-ket")
	for x, y := range base.ActionToKey {
		fmt.Printf(" Action : %v , Key : %v \n", x, sdl.GetScancodeName(y))
	}

	fmt.Println("Key-Action")
	for x, y := range base.KeyToAction {
		fmt.Printf(" Action : %v , Key : %v \n", y, sdl.GetScancodeName(x))
	}
}

func initializeControls() Inputs {
	base := Inputs{}
	base.KeyToAction = map[sdl.Scancode]Core.BUTTON{
		sdl.SCANCODE_LSHIFT: Core.ButtonSelect,
		sdl.SCANCODE_RETURN: Core.ButtonStart,
		sdl.SCANCODE_UP:     Core.ButtonUp,
		sdl.SCANCODE_DOWN:   Core.ButtonDown,
		sdl.SCANCODE_RIGHT:  Core.ButtonRight,
		sdl.SCANCODE_LEFT:   Core.ButtonLeft,
		sdl.SCANCODE_Z:      Core.ButtonA,
		sdl.SCANCODE_X:      Core.ButtonB,
	}

	base.ActionToKey = make(map[Core.BUTTON]sdl.Scancode)
	for key, action := range base.KeyToAction {
		base.ActionToKey[action] = key
	}

	return base
}

func intializeTurboControls() Inputs {
	base := Inputs{}
	base.KeyToAction = map[sdl.Scancode]Core.BUTTON{

		sdl.SCANCODE_Q: Core.ButtonA,
		sdl.SCANCODE_W: Core.ButtonB,
	}

	base.ActionToKey = make(map[Core.BUTTON]sdl.Scancode)
	for key, action := range base.KeyToAction {
		base.ActionToKey[action] = key
	}

	return base
}

func (console *game) beginTurbo(button Core.BUTTON) {

	console.TurboMux.Lock()
	defer console.TurboMux.Unlock()

	if console.TurboState[button] != nil {
		return
	}

	tk := time.NewTicker(40 * time.Millisecond)
	defer tk.Stop()

	done := make(chan bool)
	console.TurboState[button] = done

	isPressed := false

	for {
		select {
		case <-done:
			return
		case <-tk.C:
			isPressed = !isPressed
			console.core.Player1.UpdateBtnBool(button, isPressed)
		}
	}
}

func (console *game) stopTurbo(button Core.BUTTON) {

	console.TurboMux.Lock()
	defer console.TurboMux.Unlock()

	if console.TurboState[button] == nil {
		return
	}

	done := console.TurboState[button]
	done <- true

	console.TurboState[button] = nil
}
