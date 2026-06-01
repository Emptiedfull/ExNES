package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	fmt.Println("running")

	run_server()

}

type Debugger struct {
	console     *console
	Disassembly map[uint16]AssemblyLine
}

var debugConsole Debugger

func run_server() {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", startDebugger)
	mux.HandleFunc("/cpu/state", getCpuState)
	mux.HandleFunc("/cpu/lookAhead", getLookAhead)
	mux.HandleFunc("/cpu/step", runCycle)
	mux.HandleFunc("/disassembly", getDissambly)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("error running server")
	}
}

func startDebugger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	c := initializeConsole()
	c.loadROM("C:/Users/user/ExNES/games/dk.nes")
	c.Cpu.reset()

	fmt.Println("console ready for debug")

	debugConsole.console = c
	debugConsole.Disassembly = make(map[uint16]AssemblyLine)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Success")
}

func getLookAhead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	sizeStr := r.URL.Query().Get("size")
	if sizeStr == "" {
		sizeStr = "1"
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		http.Error(w, "invalid param size", http.StatusBadRequest)
		fmt.Printf("bad request: %v", err)
		return
	}

	lines := debugConsole.console.Cpu.lookAhead(debugConsole.console.Cpu.PC, size)

	err = json.NewEncoder(w).Encode(lines)
	if err != nil {
		http.Error(w, "unable to look ahead", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
}

func getDissambly(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err := json.NewEncoder(w).Encode(debugConsole.Disassembly)
	if err != nil {
		http.Error(w, "unable to parse disaambley json", http.StatusInternalServerError)
		fmt.Println("parsing json error:", err)
		return
	}
}

func runCycle(w http.ResponseWriter, r *http.Request) {

	stepsStr := r.URL.Query().Get("steps")
	if stepsStr == "" {
		stepsStr = "1"
	}
	fmt.Println(stepsStr)
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		http.Error(w, "invalid steps parameter", http.StatusBadRequest)
		return
	}

	debugConsole.console.stepCycles(steps)
	fmt.Println("ticking")

	fmt.Fprintf(w, "cycles ran succesfully %v", debugConsole.console.Cpu.PC)
}

func getCpuState(w http.ResponseWriter, r *http.Request) {
	state := debugConsole.console.Cpu.GetSate()

	err := json.NewEncoder(w).Encode(state)

	if err != nil {
		http.Error(w, "unable to get state", http.StatusInternalServerError)
		return
	}
}
