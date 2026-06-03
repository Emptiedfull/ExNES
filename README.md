# ExNES: Custom Nes Console

This is a custom built replica(more like modern copy) of the orignal nintendo console NES, including a custom made emulator in golang running on a raspberry pi along with two custom built wireless controllers and custom made cases for the components.

## Emulator Details

The emulator files reside in [The following directory](emulatorCore)

The emulator is based on the orignal nes and has the following features currently:
- Cycle implementation of the 6502 instruction set including illegal opcodes tested with SST's
- Scanline accurate pixel processing unit for static screens
- Fully implemented joypad with easy controller hooks
- Html api endpoints for various debugging functions
- NRAM support

Features planned:
- Mapper support
- Loopy registers and scrolling
- Active live debugger with breakpoint support
- Cpu optimisations for common instructions

### Structure 

```text
├── emulatorCore/
│   ├── debug_server/        # Fastapi server for console debug tools
|         ├── static/        # Contains static files to be deployed on server
|         └── main.py        # Entry Point for Fastapi App    
│   ├── logs/                # General storage for various logs 
│   └── controls.go          # JoyPad Defination and functions
│   └── cpu.go               # Cpu struct defination and functions
|   └── disassembler.go      # Utilities for cpu debugging
|   └── display.go           # Entrypoint for debug and main screen writing functions
|   └── helpers.go           # Shorthand functions for bit operaions
|   └── Main.go              # Main entrypoint
|   └── Opcodes.go           # Complete opcode map for the 6502 instruction set
|   └── ppu.go               # Main functions for the pixel processing unit
|   └── ppuDebug.go          # Utilities for debugging the ppu
|   └── server.go            # Server entrypoint for the debug server
```
** The emulator core is currently incomplete and hence does not contain running instructions and documentation**

## Controller Details 

Custom designed controller with an esp32 microcontroller with a custom designed case

Features in place : 
- Esp32 devboard microcontroller with network enabled
- Tactile push buttons for in-game inputs
- Extra buttons for debugging and QOL support
- RGB indicator for connection status
- Rechargable 3.7V 400mAh LiPo Battery with 5V boost up 
- Exposed reset button connected directly to devboard
- Custom designed 2 layer pcb with ergonomic labels

The pcb design is complete and [Build Files](controller_kicad) and [Fabrication Files](Fab) are located here

Click here to view the full [Hardware Bill of Materials (BOM)](./controller_kicad/bom.csv).

<img width="912" height="459" alt="image" src="https://github.com/user-attachments/assets/fe6b36c2-69eb-4259-8702-3a4fb2e2a532" />

<img width="737" height="437" alt="image" src="https://github.com/user-attachments/assets/c6e2ce7f-530e-4cd8-9049-a15ef333ad4d" />

The case will also be custom designed but is not currently complete

## Console Details 

The console is currently undeveloped but it has the following features planned
- Raspberry pi 2 zero brain for running emulator core
- Onboard buttons for quick actions such as pairing and reseting 
- Visual indicator for controller connections 
- Hdmi output compatiblity
- Custom designed enclosure with heatsink support 

Requirements for the console may change as development begins
