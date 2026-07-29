# Features

This page is a comphrensive list of all available features (hopefully)

## Console Controls 

![hehe](./icons/legend.png)

### Power button (1)

Controls the physical state of being **powered**. It also is the start of the console. Turning it on and off stimulates a **power cycle** like on the orignal NES.

## Reset (2)

Works like the **reset** button on the physical NES. Only works in audio clocked mode. May cause data corruption in some games (this is not a glitch, this is a artifact of accurate emulation)

## Play/Pause (3)

Can be used to control the state of emulation. The pausing action is not instant and may cause **5-10 frames** of lag. Only works in audio clocked mode.
 

### Power led and ROM Led (4 and 5)

This represent the state of the console at a given point. Note any flickering effects are purely **cosmetic** and do not affect 
usage.

If your rom has been loaded but the power led remains blank, use the **power button** to turn on the console.


### Controller ports (6 and 7)

Uh my bad they aren't really implemented yet but they would allow you to use external control devices and connect and disconnect input sources. Have a happy cat on me instead. 

### Catridge slot (8)

Press on the **slot** to open up the rom selection menu. Using it mid game can be used to hotswap games which can be used for some exploits in particular games, it can also however cause data corruption. Restart the tab to fix this.


## Television controls 

![altTextisCringeSorry](./static/icons/tv.png)

### Speed dial(1)

It can be used to control the emulation speed to slow or speed it up, by defualt the increment are in orders of 50%. Only works in audio clocked mode.

### Volume dial(2)

Pretty obvious, controls the audio. Only works in audio clocking mode.

### Tv Screen

Clicking on the screen changes the display mode to fullscreen

## Miscellanous 

### Timing control

Probably the most drastic control option you have the liberty to change. Check out [here](./docs/timing.md) for more info

### Updating controls 

The on screen controller can be used to change up the controls, simply press the button and then rebind it with any keyboard key (or hopefully external controller but no promises). The onscreen controller can also be used to actually control the emulator btw!

### Rewind

You found the easter egg!!!! Press the R key on the keyboard to open up the rewind overlay, you can then use the timeline to scrub into a past frame from a collection of the last 4 seconds of frames. The frame can then be loaded into for rewinding, enjoy cheating you lil bitch.

