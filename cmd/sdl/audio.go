package main

/*
#cgo pkg-config: sdl2
#include <SDL2/SDL.h>
*/
import "C"
import "unsafe"

//export audioCallback
func audioCallback(userdata unsafe.Pointer, stream *C.uchar, length C.int) {
	n := int(length) / 2
	buf := unsafe.Slice((*int16)(unsafe.Pointer(stream)), n)

	for i := range n {
		for !console.Apu.HasSample() {
			console.TickNoAudio()
		}
		buf[i] = convertAudio(console.Apu.PopSample())
	}
}

func convertAudio(sample float32) int16 {
	center := (sample * 2.0) - 1.0
	if center > 1.0 {
		center = 1.0
	} else if center < -1.0 {
		center = -1.0
	}
	return int16(center * 32767)
}
