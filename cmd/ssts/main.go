package main

import (
	"encoding/json"
	"exnes/Core"
	_ "exnes/Core"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type TestItem struct {
	Name    string      `json:"name"`
	Initial cpuState    `json:"initial"`
	Final   cpuState    `json:"final"`
	Cycles  []cycleInfo `json:"cycles"`
}

type cpuState struct {
	Pc  int      `json:"pc"`
	S   int      `json:"s"`
	X   int      `json:"x"`
	Y   int      `json:"y"`
	A   int      `json:"a"`
	P   int      `json:"p"`
	Ram [][2]int `json:"ram"`
}

type cycleInfo struct {
	Addr int
	Val  int
	Op   string
}

func (c *cycleInfo) UnmarshalJSON(data []byte) error {
	var message [3]json.RawMessage

	if err := json.Unmarshal(data, &message); err != nil {
		return err
	}

	json.Unmarshal(message[0], &c.Addr)
	json.Unmarshal(message[1], &c.Val)
	json.Unmarshal(message[2], &c.Op)

	return nil

}

func loadTestJson(filepath string) ([]TestItem, error) {

	var tests []TestItem

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &tests)
	if err != nil {
		return nil, err
	}

	return tests, nil
}

func runTest(filepath string) error {
	TestArr, err := loadTestJson(filepath)
	if err != nil {
		return nil
	}

	c := Core.Cpu{}
	c.Mem = Core.GetBus()
	c.Mem.FlatMode = true
	c.Mem.FlatMem = make([]uint8, 65536)

	c.Mem.Cpu = &c

	for _, test := range TestArr {

		loadCpuState(&c, test.Initial)

		for range test.Cycles {
			c.Tick()
		}

		err := compareCycles(c.Mem.Log, test.Cycles)
		if err != nil {
			return err
		}

		err = CompareState(&c, test.Final)
		if err != nil {
			return err
		}

		c.Mem.Log = make([]Core.CycleStep, 0)
	}
	fmt.Println("test passed:", filepath)

	return nil
}

func loadCpuState(cpu *Core.Cpu, state cpuState) {

	for _, kv := range state.Ram {
		cpu.Mem.FlatMem[uint16(kv[0])] = uint8(kv[1])
	}

	cpu.A = uint8(state.A)
	cpu.X = uint8(state.X)
	cpu.Y = uint8(state.Y)
	cpu.PC = uint16(state.Pc)
	cpu.P = uint8(state.P)
	cpu.S = uint8(state.S)
	cpu.ResetTemp()
}

func compareCycles(cpuLog []Core.CycleStep, testLog []cycleInfo) error {
	if len(cpuLog) != len(testLog) {
		return fmt.Errorf("cycles not matching, cpuLog: %v testLog: %v", cpuLog, testLog)
	}
	return nil
}

func CompareState(cpu *Core.Cpu, state cpuState) error {
	for _, kv := range state.Ram {
		val := cpu.Mem.FlatMem[kv[0]]
		if val != uint8(kv[1]) {
			return fmt.Errorf("Invalid Value at %v, wanted: %v got: %v", kv[0], kv[1], val)
		}
	}

	if cpu.A != uint8(state.A) {
		return fmt.Errorf("A register mismatch: wanted %v got %v", state.A, cpu.A)
	}

	if cpu.X != uint8(state.X) {
		return fmt.Errorf("X no good, wanted %v got %v", state.X, cpu.X)
	}

	if cpu.Y != uint8(state.Y) {
		return fmt.Errorf("Y no good,wanted %v got %v", state.Y, cpu.Y)
	}

	if cpu.PC != uint16(state.Pc) {
		return fmt.Errorf("Pc wrong, wanted %v got %v", state.Pc, cpu.PC)
	}

	if cpu.S != uint8(state.S) {
		return fmt.Errorf("Stack wrong, wanted %v got %v", state.S, cpu.S)
	}

	if cpu.P != uint8(state.P) {
		return fmt.Errorf("P bad, wanted %v got %v", state.P, cpu.P)
	}

	return nil
}

func StartTesting() {
	var wg sync.WaitGroup
	filepath.WalkDir("./nes6502/v1", func(path string, d fs.DirEntry, err error) error {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := runTest(path)
			if err != nil {
				log.Fatalf("Error on test: %v, error: %v", path, err)
			}
		}()

		// runTest(path)

		return nil
	})

	wg.Wait()
}

func main() {
	// err := runTest("/Users/test/Projects/ExNES/cmd/ssts/nes6502/v1/0a.json")
	// if err != nil {
	// 	panic(err)
	// }

	StartTesting()
}
