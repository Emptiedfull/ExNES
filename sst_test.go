package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

type codeTest struct {
	Name    string      `json:"name"`
	Initial cpustate    `json:"initial"`
	Final   cpustate    `json:"final"`
	Cycles  []cycleStep `json:"cycles"`
}

func (c *cycleStep) UnmarshalJSON(data []byte) error {
	var interfaceSlice []interface{}

	if err := json.Unmarshal(data, &interfaceSlice); err != nil {
		return err
	}

	c.Addr = uint16(interfaceSlice[0].(float64))
	c.Val = uint8(interfaceSlice[1].(float64))
	c.Mode = interfaceSlice[2].(string)

	return nil
}

type cpustate struct {
	Pc  uint16 `json:"pc"`
	S   uint8  `json:"s"`
	A   uint8  `json:"a"`
	X   uint8  `json:"x"`
	Y   uint8  `json:"y"`
	P   uint8  `json:"p"`
	Ram [][]int
}

const Dir = "C:/Users/user/ExNES/65x02/nes6502/v1/"

func loadJson(t *testing.T, filepath string) []codeTest {
	file, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("failed to read file data %v", err)
	}

	var tests []codeTest
	if err := json.Unmarshal(file, &tests); err != nil {
		t.Fatalf("unable to parse test json")
	}

	return tests
}

func getTestNames(t *testing.T) []string {
	entries, err := os.ReadDir(Dir)
	if err != nil {
		t.Fatalf("unable to read file %v", err)
	}

	var testNames []string

	for _, entry := range entries {
		testNames = append(testNames, entry.Name())
	}

	return testNames
}

func TestOpcodes(t *testing.T) {

	tests := getTestNames(t)

	for _, test := range tests[:255] {
		performTest(t, test)
	}

	fmt.Println("ALL TESTS PASSED")

}

func performTest(t *testing.T, testName string) {

	testDir := Dir + testName

	tests := loadJson(t, testDir)

	c := cpu{
		mem: &TestBus{},
	}

	for _, test := range tests {
		c.LoadCpuState(test.Initial)

		cyclecount := len(test.Cycles)

		for range cyclecount - 1 {
			c.tick()
		}

		final, ram := c.GetCpuState()
		_, err := compareStates(ram, final, test.Final)
		if err != nil {
			t.Errorf("Error occured at test: %v, %s \n cycles: %v \n performed: %v", test.Name, err.Error(), test.Cycles, ram.GetHistory())
			return
		}

		c.mem.ClearHistory()
		c.isJamming = false

	}

	t.Log(testName, "PASSED")
}

func (TestBus *TestBus) GetHistory() []cycleStep {
	return TestBus.History
}

func (TestBus *TestBus) ClearHistory() {
	TestBus.History = []cycleStep{}
}

func compareStates(ram typebus, got, expected cpustate) (bool, error) {

	if got.Pc != expected.Pc {
		return false, fmt.Errorf("Pc mismatch, got: %v expected %v", got.Pc, expected.Pc)
	}

	if got.A != expected.A {
		return false, fmt.Errorf("Register A mismatch got: %v expected %v", got.A, expected.A)
	}

	if got.X != expected.X {
		return false, fmt.Errorf("Register X mismatch, got %v expected %v", got.X, expected.X)
	}

	if got.Y != expected.Y {
		return false, fmt.Errorf("Register Y mismatch, got %v expected %v", got.Y, expected.Y)
	}

	if got.S != expected.S {
		return false, fmt.Errorf("Stack pointer mismatch, got %v expected %v", got.S, expected.S)
	}

	if got.P != expected.P {
		return false, fmt.Errorf("Flag mismatch, got %v expected %v", got.P, expected.P)
	}

	for _, memCell := range expected.Ram {
		addr := memCell[0]
		expectedVal := memCell[1]

		if ram.Get(uint16(addr)) != uint8(expectedVal) {
			return false, fmt.Errorf("Memory mismatch at %b, expected %b got %b", addr, expectedVal, ram.Get(uint16(addr)))
		}

	}

	return true, nil
}

type TestBus struct {
	RAM     [65536]byte
	History []cycleStep
}

func (T *TestBus) Read(addr uint16) uint8 {
	val := T.RAM[addr]
	T.History = append(T.History, cycleStep{Addr: addr, Val: val, Mode: "read"})
	return val
}

func (T *TestBus) Write(addr uint16, val uint8) {
	T.Set(addr, val)
	T.History = append(T.History, cycleStep{Addr: addr, Val: val, Mode: "write"})
}

func (c *cpu) LoadCpuState(state cpustate) {
	c.PC = state.Pc
	c.S = state.S
	c.A = state.A
	c.X = state.X
	c.Y = state.Y
	c.P = state.P

	for _, cell := range state.Ram {
		addr := uint16(cell[0])
		val := uint8(cell[1])
		c.mem.Set(addr, val)
	}
}

func (c *cpu) GetCpuState() (cpustate, typebus) {
	state := cpustate{}

	state.A = c.A
	state.X = c.X
	state.Y = c.Y
	state.Pc = c.PC
	state.P = c.P
	state.S = c.S
	return state, c.mem
}

func (T *TestBus) Set(addr uint16, val uint8) {
	T.RAM[addr] = val
}

func (T *TestBus) Get(addr uint16) uint8 {
	return T.RAM[addr]
}
