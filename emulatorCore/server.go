package main

import (
	"encoding/json"
	"exnes/Core"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	fmt.Println("running")

	setUpListener()

}

var debugConsole Core.Debugger

func run_server() {
	mux := http.NewServeMux()

	mux.HandleFunc("/Debugger/reset", startDebugger)
	mux.HandleFunc("/cpu/state", getCpuState)
	mux.HandleFunc("/cpu/lookAhead", getLookAhead)
	mux.HandleFunc("/cpu/step", runCycle)
	mux.HandleFunc("/disassembly", getDissambly)
	mux.HandleFunc("/screen/get/Debug", getDebugScreen)
	mux.HandleFunc("/ppu/debugCHR", runChrViewer)
	mux.HandleFunc("/ppu/debugNameTable", runNameTableViewer)
	mux.HandleFunc("/screen", getScreenIMM)
	// mux.HandleFunc("/controls/update", updateControls)
	mux.HandleFunc("/screen/socket", acceptScreenConn)
	mux.HandleFunc("/start/console", startConsole)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("error running server")
	}
}

func startConsole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if debugConsole.Console == nil {
		http.Error(w, "console not initialized", http.StatusInternalServerError)
		return
	}

	go debugConsole.Console.StartConsoleCycle()

	w.WriteHeader(http.StatusOK)
}

func updateControls(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var state Core.ControlState

	err := json.NewDecoder(r.Body).Decode(&state)

	if err != nil {
		http.Error(w, "Unable to update controls", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	if debugConsole.Console == nil {
		http.Error(w, "Console not started", http.StatusInternalServerError)
		return
	}

	debugConsole.Console.JoyPad.UpdateState(state)

	w.WriteHeader(http.StatusOK)
}

func getScreenIMM(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/octet-stream")

	if debugConsole.Console == nil {
		http.Error(w, "unable to get screen buffer", http.StatusBadRequest)
		return
	}

	w.Write(debugConsole.Console.Ppu.GetScreenBuffer())

}

func runChrViewer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	debugConsole.Console.DebugChrRom()
	fmt.Fprint(w, "Success")
}

func runNameTableViewer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	debugConsole.Console.DebugNameTable()
	fmt.Fprintf(w, "Success")

}

func getDebugScreen(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/octet-stream")

	if debugConsole.Console == nil {
		http.Error(w, "unable to get screen buffer", http.StatusBadRequest)
		return
	}

	w.Write(debugConsole.Console.Ppu.DebugBuffer)
}

func startDebugger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	c := Core.InitializeConsole()
	c.LoadROM("C:/Users/user/ExNES/emulatorCore/games/dk3.nes")
	c.Cpu.Reset()

	fmt.Println("console ready for debug")

	debugConsole.Console = c
	debugConsole.Disassembly = make(map[uint16]Core.AssemblyLine)
	debugConsole.Console.ScreenChannel = make(chan Core.ScreenInfo, 100)

	c.Ppu.DebugBuffer = make([]uint8, 512*64)

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

	lines := debugConsole.LookAhead(debugConsole.Console.Cpu.PC, size)

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	stepsStr := r.URL.Query().Get("steps")
	if stepsStr == "" {
		stepsStr = "1"
	}
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		http.Error(w, "invalid steps parameter", http.StatusBadRequest)
		return
	}

	debugConsole.StepCycles(steps)

	fmt.Fprintf(w, "cycles ran succesfully %x", debugConsole.Console.Cpu.PC)

}

func getCpuState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if debugConsole.Console == nil {
		http.Error(w, "console not started", http.StatusInternalServerError)
		return
	}

	state := debugConsole.Console.Cpu.GetSate()

	err := json.NewEncoder(w).Encode(state)

	if err != nil {
		http.Error(w, "unable to get state", http.StatusInternalServerError)
		return
	}
}
