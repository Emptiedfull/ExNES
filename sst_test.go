package main_test

import (
	"encoding/json"
	"os"
	"testing"
)

type codeTest struct {
	Name    string      `json:"name"`
	Initial cpustate    `json:"initial"`
	Final   cpustate    `json:"final"`
	Cycles  []cycleStep `json:"cycles"`
}

type cycleStep struct {
	Addr uint16
	Val  uint8
	Mode string
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

func loadJson(t *testing.T, filepath string) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("failed to read file data %v", err)
	}

	var tests []codeTest
	if err := json.Unmarshal(file, &tests); err != nil {
		t.Fatalf("unable to parse test json")
	}

	t.Log(tests[0].Cycles)
}

func TestOpcodes(t *testing.T) {
	loadJson(t, "C:/Users/user/ExNES/ProcessorTests/6502/v1/00.json")
}
