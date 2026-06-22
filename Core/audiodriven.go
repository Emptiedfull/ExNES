package Core

import (
	"fmt"
	"math"
	"sync/atomic"
	"time"
)

func (a *APU) MalgoAdapter(output []byte, input []byte, framecount uint32) {
	a.DriveSamples(output, framecount)
}

func (a *APU) DriveSamples(output []byte, samplesNeeded uint32) {
	for i := range samplesNeeded {
		for !a.HasSample() {
			a.Console.TickNoAudio()
		}

		sample := a.PopSample()

		bits := math.Float32bits(sample)
		output[i*4] = byte(bits)
		output[i*4+1] = byte(bits >> 8)
		output[i*4+2] = byte(bits >> 16)
		output[i*4+3] = byte(bits >> 24)
	}
}
func (c *Console) TickNoAudio() {
	if c.Cpu.Stall > 0 {
		c.Cpu.Stall--
		c.Cpu.TotalCycles++
	} else {
		c.Cpu.tick()
	}

	for range 3 {
		c.Ppu.step()
	}

	c.Apu.tick()
	if c.Apu.IRGPending || c.Apu.Dmc.IRGPending {
		c.Cpu.triggerIRQ()
	}

	// c.RunDisplayUpdates()

}

type AudioDebug struct {
	ReadCalls      atomic.Int64
	BytesPulled    atomic.Int64
	SamplesCreated atomic.Int64
	Underruns      atomic.Int64
}

var AudioStats AudioDebug

func (a *APU) LogAudioStats() {

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		calls := AudioStats.ReadCalls.Swap(0)
		bytes := AudioStats.BytesPulled.Swap(0)
		// samples := AudioStats.SamplesCreated.Swap(0)
		underruns := AudioStats.Underruns.Swap(0)

		samplesPerCall := int64(0)
		if calls > 0 {
			samplesPerCall = (bytes / 4) / calls
		}

		fmt.Printf(
			"[APU] reads/s: %d | bytes/s: %d | samples pulled/s: %d (want ~44100) | samples/read: %d | underruns/s: %d | buffer backlog: %d\n",
			calls, bytes, bytes/4, samplesPerCall, underruns, len(a.SampleBuffer),
		)
	}
}
