package Core

import (
	"math"
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

	if a.Console.Ppu.ScreenChanged {
		a.Console.RunDisplayUpdates()

		if a.Console.Ppu.Frame%20 == 0 {

			a.Console.TakeSnapshot()
		}

	}

}
func (c *Console) TickNoAudio() {
	if c.Cpu.Stall > 0 {
		c.Cpu.Stall--
		c.Cpu.TotalCycles++
	} else {
		c.Cpu.Tick()
	}

	for range 3 {
		c.Ppu.step()
	}

	c.Apu.tick()
	if c.Apu.IRGPending || c.Apu.Dmc.IRGPending || c.Cpu.currentstep == 0 {
		c.Cpu.triggerIRQ()
	}

}
