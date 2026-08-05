import { wait } from "./joypad.js"
import { createModal } from "./modal.js"
import { createTilesFromSnapshots } from "./rewind.js"

const worker = new Worker(new URL('./emuWorker.js', import.meta.url))
let fBytes = null

const inputBuf = new SharedArrayBuffer(4)
const inputState = new Int32Array(inputBuf)

const cmdMap = {
    "STOP":1,
    "RESET":2,
}

const canvas = document.getElementById("screen")
const ctx = canvas.getContext("2d")
const imageData = ctx.createImageData(256, 240)

const speedBuf = new SharedArrayBuffer(4)
const speedNum = new Int32Array(speedBuf)

const NES_FPS = 60.0948
const FRAME_MS = 1000 / NES_FPS

let frameSig = null
let rafMode = false
let rafId = null
let lastT = 0
let  acc = 0

let audioBufS = null
let AudioControl = null

let gameControl = null

Atomics.store(speedNum,0,1000)

export const state = {
    romRunning: false,
    runMode:0, //0-aduio run 1 - raf run
}

let audioCtx = null
let gain = null

let startTime = null

window.addEventListener("keydown",async (e)=>{
    if (e.code == "KeyR"){
       
        await getSnapList()
        
       
       
    } else if (e.code == "KeyT"){
        await PauseGame()
        
        worker.postMessage({type:"reset"})

       
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

    // if (audioCtx == null) {
    //     return
    // }
    // await audioCtx.suspend()

    

    Atomics.store(gameControl,0,2)
    Atomics.notify(gameControl,0)

}

export const ResumeGame = async ()=>{
    
    if (audioCtx == null){
        return
    }
    await audioCtx.resume()

    Atomics.store(gameControl,0,cmdMap.STOP)
    Atomics.notify(gameControl,0)

    worker.postMessage({"type":"pump"})
}

export const getSnapList = async ()=>{
    await PauseGame()

    worker.postMessage({type:"getsnap"})
    
}

export const loadSnap = async(index)=>{
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
            AudioControl = new Int32Array(audioBufS,data.SIZE*4,3)
            gameControl = new Int32Array(data.gameBuf)
            frameSig = new Int32Array(data.frameSigBuf)
            await setUpAudio(data.audioBufS, data.SIZE)
            break
        case "wasm":
            createModal("Console Ready!!","Press the power button to start the console")
            worker.postMessage({ type: "init", inputBuf,speedBuf })
            break
        case 'snaps':
            await createTilesFromSnapshots(data.snaps)
            break
        case 'start':
            console.log("starting/restarting")
            await startConsole()
            break

        case 'frameUp':
            imageData.data.set(fBytes)
            ctx.putImageData(imageData, 0, 0)
            break
    }
}

const startConsole = async ()=>{
    await ResumeGame()
    if (state.romRunning){
        return
    }

     state.romRunning = true
    console.log("starting the fucking game")
    if (state.runMode == 0){
        await startAudioBuf()
    }else{
        await startRaf()
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

export const switchMode = (mode)=>{
    if (mode == 0){
        state.runMode= 0
        startAudioBuf()
    }else{
        state.runMode = 1
        startRaf()
    }
}

export const loadRom = async (game) => {
  
    // Atomics.store(control,2,1)
    // Atomics.notify(control,2)
    await PauseGame()

    // audioCtx.resume()
    console.log("loading game")

    worker.postMessage({ type: 'loadRom', rom: game })
    
    // await wait(1000)
    
   
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

export const startAudioBuf = async  ()=>{

    if (!state.romRunning) return
   stopRaF()

    if (audioCtx) {
        await audioCtx.resume()
    }

    await ResumeGame()

    worker.postMessage({type:"pump"})
    
}

export const startRaf = ()=>{
    if (!state.romRunning) return
    if (!frameSig) return

    if (gameControl) {
        Atomics.store(gameControl,0,cmdMap.STOP)
        Atomics.notify(gameControl,0)
    }

    if (audioCtx) audioCtx.suspend()

    Atomics.store(frameSig, 2, 0)  
    Atomics.store(frameSig, 1, Atomics.load(frameSig, 0)) 
    rafMode = true
    lastT = performance.now()
    acc = 0

     worker.postMessage({type:"startRaF"})

    rafId = requestAnimationFrame(rafLoop)

}



const rafLoop = (now)=>{
    if (!rafMode) return

    acc += now - lastT
    lastT = now 

    if (acc > 250) acc = 250

    while (acc > FRAME_MS){
        Atomics.add(frameSig,0,1)
        Atomics.notify(frameSig,0)

        acc -= FRAME_MS
    }

    rafId = requestAnimationFrame(rafLoop)
}

export const stopRaF = ()=>{
    if (rafId) cancelAnimationFrame(rafId)

    rafId = null
    rafMode = false
    if (frameSig){
        Atomics.store(frameSig,2,1)
        Atomics.add(frameSig,0,1)
        Atomics.notify(frameSig,0)
    }
}