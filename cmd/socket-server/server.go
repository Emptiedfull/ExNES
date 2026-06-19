package main

import (
	"encoding/json"
	"exnes/Core"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
)

func main() {

	fmt.Println("running")

	run_server()

}

var debugConsole Core.Debugger

func run_server() {
	mux := http.NewServeMux()

	mux.HandleFunc("/cpu/state", getCpuState)
	mux.HandleFunc("/cpu/lookAhead", getLookAhead)
	mux.HandleFunc("/disassembly", getDissambly)
	mux.HandleFunc("/screen/get/Debug", getDebugScreen)

	mux.HandleFunc("/controls/update", updateControls)
	mux.HandleFunc("/screen/socket", acceptScreenConn)

	mux.HandleFunc("/Debugger/reset", startDebugger)
	mux.HandleFunc("/Debugger/getConsoleStatus", getConsoleStatus)

	mux.HandleFunc("/console/start", startConsole)
	mux.HandleFunc("/console/pause", pauseConsole)
	mux.HandleFunc("/console/unpause", unpauseConsole)
	mux.HandleFunc("/console/getExecStatus", getExecStatus)

	mux.HandleFunc("/run/cycle", runCycle)
	mux.HandleFunc("/run/frame", runFrame)
	mux.HandleFunc("/run/frame30", run30frame)

	mux.HandleFunc("/snapshots/fetch", fetchSnapshots)
	mux.HandleFunc("/snapshots/load", loadSnapshot)

	mux.HandleFunc("/ppu/debugCHR", runChrViewer)
	mux.HandleFunc("/ppu/debugNameTable", runNameTableViewer)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("error running server")
	}
}

func getConsoleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := json.NewEncoder(w).Encode(debugConsole.Console != nil)
	if err != nil {
		http.Error(w, "error getting console status", http.StatusInternalServerError)
		return
	}
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

func getExecStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if debugConsole.Console == nil {
		http.Error(w, "console not initialized", http.StatusInternalServerError)
		return
	}

	err := json.NewEncoder(w).Encode(debugConsole.Console.Paused)
	if err != nil {
		http.Error(w, "unable to encode pause status-very weird", http.StatusInternalServerError)
		fmt.Println("something really bad happened encoding a bool failing is insane")
		return
	}

}

func runSingleCycle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	debugConsole.DebugTick()
	debugConsole.Console.RunDisplayUpdates()

	w.WriteHeader(http.StatusOK)
}

func run30frame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	for range 29781 * 30 {
		debugConsole.DebugTick()
	}

	debugConsole.Console.RunDisplayUpdates()

	w.WriteHeader(http.StatusOK)
}

func runFrame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	debugConsole.RunDebugFrame()
	w.WriteHeader(http.StatusOK)
}

func startConsole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if debugConsole.Console == nil {
		http.Error(w, "console not initialized", http.StatusInternalServerError)
		return
	}

	// go debugConsole.StartDebugConsole()
	setUpAudioDriver()
	// go debugConsole.Console.Apu.LogAudioStats()

	w.WriteHeader(http.StatusOK)
}

func pauseConsole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if debugConsole.Console == nil {
		http.Error(w, "console not initialized", http.StatusInternalServerError)
		return
	}

	debugConsole.Console.Pause()

	w.WriteHeader(http.StatusOK)

}

func unpauseConsole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if debugConsole.Console == nil {
		http.Error(w, "console not initialized", http.StatusInternalServerError)
		return
	}

	debugConsole.Console.UnPause()

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

	debugConsole.Console.Player1.UpdateState(state)

	w.WriteHeader(http.StatusOK)
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

func quickStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	debugConsole.Console = Core.Quickstart("C:/Users/user/ExNES/emulatorCore/games/Mapper2/contra.nes")

	go HandleScreenUpdates()

	w.WriteHeader(http.StatusOK)

}

func startDebugger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	debugConsole.Console = Core.Quickstart("C:/Users/user/ExNES/games/NROM/mario.nes")
	// c.LoadROM("C:/Users/user/ExNES/emulatorCore/games/Mapper2/contra.nes")

	//c.LoadROM("C:/Users/user/ExNES/emulatorCore/games/Mapper1/ff.nes")
	// c.LoadROM("C:/Users/user/ExNES/emulatorCore/test_roms/ppu/vbl.nes")
	// c.LoadROM("C:/Users/user/ExNES/emulatorCore/games/Mapper4/mario3.nes")

	debugConsole.Disassembly = make(map[uint16]Core.AssemblyLine)

	debugConsole.Console.Ppu.DebugBuffer = make([]uint8, 512*64)

	go HandleScreenUpdates()

	w.WriteHeader(http.StatusOK)

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
