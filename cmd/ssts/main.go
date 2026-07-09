package main

import (
	"encoding/json"
	"exnes/Core"
	_ "exnes/Core"
	"os"
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

}

func main() {

}
