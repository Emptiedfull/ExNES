const worker = new Worker("/static/scripts/extra/emuWorker.js")
let fBytes = null

const inputBuf = new SharedArrayBuffer(4)
const inputState = new Int32Array(inputBuf)

let romRunning = false


const setUpAudio = async (audioBufS,SIZE)=>{
    const audioCtx = new AudioContext({sampleRate:44100})

    await audioCtx.audioWorklet.addModule("/static/scripts/extra/driverWorklet.js")

    const node = new AudioWorkletNode(audioCtx,'apu-proc',{
        outputChannelCount : [1],
    })

    node.port.postMessage({audioBufS,SIZE})
    node.connect(audioCtx.destination)
}


const initConsole = async ()=>{
    worker.postMessage({type:"init",inputBuf})
}

worker.onmessage = async ({data})=>{
   
    switch (data.type){
        case 'init':
            fBytes = new Uint8Array(data.FBuf)
            await setUpAudio(data.audioBufS,data.SIZE)
            break
        case 'frameUp':
            imageData.data.set(fBytes)
            ctx.putImageData(imageData, 0, 0)
            break
    }
}

const ControlMap = {
    "joypad-A":0,
    "joypad-B":1,
    "joypad-select":2,
    "joypad-start":3,
    "dpad-up":4,
    "dpad-down":5,
    "dpad-left":6,
    "dpad-right":7,
}

const loadRom = (game) => {
    worker.postMessage({ type: 'loadRom', rom: game })
}

const UpdatePress = (btn)=>{
   
    if (romRunning){
        mask = ControlMap[btn]
        Atomics.or(inputState,0,1 << mask)
    }
}

const UpdateRelease = (btn) =>{
    
    if (romRunning){
        mask = ControlMap[btn]
        
        Atomics.and(inputState,0,~(1 << mask))
    }
}