import { wait } from "./joypad.js"
import { createModal } from "./modal.js"

const worker = new Worker(new URL('./emuWorker.js', import.meta.url))
let fBytes = null

const inputBuf = new SharedArrayBuffer(4)
const inputState = new Int32Array(inputBuf)

const canvas = document.getElementById("screen")
const ctx = canvas.getContext("2d")
const imageData = ctx.createImageData(256, 240)

const speedBuf = new SharedArrayBuffer(4)
const speedNum = new Int32Array(speedBuf)

let audioBufS = null
let control = null

Atomics.store(speedNum,0,1000)

export const state = {
    romRunning: false,
}

let audioCtx = null
let gain = null

let startTime = null

window.addEventListener("keydown",async (e)=>{
    if (e.code == "KeyL"){
        startTime = performance.now()
        await getSnapList()
    } else if (e.code == "KeyP"){
        console.log("hello")
        await loadSnap(0)
    }
})

const setUpAudio = async (audioBufS, SIZE) => {

    if (audioCtx !== null && gain !== null){
        return 
    }

    audioCtx = new AudioContext({ sampleRate: 44100 })

    await audioCtx.audioWorklet.addModule(new URL("./driverWorklet.js", import.meta.url))

    const node = new AudioWorkletNode(audioCtx, 'apu-proc', {
        outputChannelCount: [1],
    })

    gain = audioCtx.createGain()

    node.port.postMessage({ audioBufS, SIZE })


    node.connect(gain)
    gain.connect(audioCtx.destination)

}

export const PauseGame = async ()=>{

    if (audioCtx == null) {
        return
    }
    await audioCtx.suspend()

    Atomics.store(control,2,2)
    Atomics.notify(control,2)
}

export const ResumeGame = async ()=>{
    console.log("playing")
    if (audioCtx == null){
        return
    }
    await audioCtx.resume()

    Atomics.store(control,2,1)
    Atomics.notify(control,2)

    worker.postMessage({"type":"pump"})
}

export const getSnapList = async ()=>{
    await PauseGame()

    worker.postMessage({type:"getsnap"})
    
}

export const loadSnap = async(index)=>{
     await PauseGame()
    worker.postMessage({type:"loadSnapshot",index})

    await wait(100)
    await ResumeGame()
}

export const updateVolume = (intensity) =>{
    if (gain !== null){
           gain.gain.setValueAtTime(gain.gain.value, audioCtx.currentTime);
           gain.gain.linearRampToValueAtTime(intensity, audioCtx.currentTime + 1)
    }
}

export const initConsole = ()=>{
    worker.postMessage({type:"init",speedBuf,inputBuf})
}



worker.onmessage = async ({ data }) => {

    switch (data.type) {
        case 'init':
            fBytes = new Uint8Array(data.FBuf)
            audioBufS = data.audioBufS
            control = new Int32Array(audioBufS,data.SIZE*4,3)
            await setUpAudio(data.audioBufS, data.SIZE)
            break
        case "wasm":
            createModal("Console Ready!!","Press the power button to start the console")
            worker.postMessage({ type: "init", inputBuf,speedBuf })
            break
        case 'snaps':
            console.log(data)
            console.log(performance.now() - startTime)
            await ResumeGame()
            break
        case 'frameUp':
            imageData.data.set(fBytes)
            ctx.putImageData(imageData, 0, 0)
            break
    }
}

const ControlMap = {
    "joypad-A": 0,
    "joypad-B": 1,
    "joypad-select": 2,
    "joypad-start": 3,
    "dpad-up": 4,
    "dpad-down": 5,
    "dpad-left": 6,
    "dpad-right": 7,
}

export const loadRom = (game) => {
    Atomics.store(control,2,1)
    Atomics.notify(control,2)

    audioCtx.resume()

    worker.postMessage({ type: 'loadRom', rom: game })
}

export const UpdateSpeed = (speed) =>{
    Atomics.store(speedNum,0,speed)
}

export const UpdatePress = (btn) => {

    Atomics.or(inputState, 0, 1 << ControlMap[btn])

}

export const UpdateRelease = (btn) => {

    Atomics.and(inputState, 0, ~(1 << ControlMap[btn]))

}