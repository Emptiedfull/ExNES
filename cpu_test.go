package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func LoadNesTest(c *cpu) {

	data, err := os.ReadFile("nestest.nes")
	if err != nil {
		log.Fatal("failed to read file", err)
	}

	prgStart := 16
	prgEnd := 16 + 16384

	c.mem.external = data[prgStart:prgEnd]

	c.PC = 0xC000
	c.S = 0xFD
	c.P = 0x24
	c.A = 0
	c.X = 0
	c.Y = 0
	c.currentstep = 0
	c.totalCycles = 7

	fmt.Println("Test loaded")
}

func (c *cpu) logState() string {
	opC := c.mem.Read(c.PC)
	b1 := c.mem.Read(c.PC + 1)
	b2 := c.mem.Read(c.PC + 2)

	Name := FetchTable[opC].Name

	return fmt.Sprintf("%04X  %02X %02X %02X  %-3s A:%02X X:%02X Y:%02X P:%02X SP:%02X CYC:%d",
		c.PC, opC, b1, b2, Name, c.A, c.X, c.Y, c.P, c.S, c.totalCycles)
}

func TestNesTest(t *testing.T) {
	cpu := &cpu{}
	cpu.mem = &bus{}

	LoadNesTest(cpu)

	masterFile, err := os.Open("masterlog.txt")
	if err != nil {
		t.Fatalf("Could not open file: %v", err)
	}

	outFile, err := os.Create("my_emu_nestest.log")
	if err != nil {
		t.Fatalf("Failed to create output log file: %v", err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	defer masterFile.Close()
	Scanner := bufio.NewScanner(masterFile)

	for i := 0; i < 8000; i++ {
		if cpu.currentstep == 0 {
			if Scanner.Scan() {
				expectedline := Scanner.Text()
				myLine := cpu.logState()

				_, _ = writer.WriteString(myLine + "\n")

				if myLine[:4] != expectedline[:4] {
					t.Errorf("Instruction %d: PC Mismatch!\nGot:  %s\nWant: %s", i, myLine, expectedline)
					return
				}

				myCycIdx := strings.LastIndex(myLine, "CYC:")
				expCycIdx := strings.LastIndex(expectedline, "CYC:")

				if myCycIdx != -1 && expCycIdx != -1 {
					myCycles := strings.TrimSpace(myLine[myCycIdx+4:])
					expCycles := strings.TrimSpace(expectedline[expCycIdx+4:])

					if myCycles != expCycles {
						t.Errorf("Instruction %d: Cycle Mismatch!\nGot:  %s\nWant: %s", i, myLine, expectedline)
						return
					}
				}

			}
		}

		cpu.tick()

		for cpu.currentstep != 0 {
			cpu.tick()
		}
	}

}
