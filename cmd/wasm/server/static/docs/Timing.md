# Timing Modes

ExNES supports a two different ways of timing the emulator clock. This should not concern most people but it can be the source of your problem. 

## Which mode to use? 

The default mode at start up is Audio clocked, this is because it has the most extensive feature support. However you might find that using audio clocked makes the game jittery try to switch to rAF Clocked.

## How to switch modes? 

Gng, its the most akward button on the screen, just use the slider beneath the hero text.

## Comparison 

|     | Audio Clocked | rAF Clocked |
|-------|-------|-------|
| Compatiblity    | Medium    | High     |
| Frame rate   | Frame perfect     |  ~60 fps     |
|Audio | Supported | Not Supported | 


