let wasm_up = false
importScripts('/static/wasm_exec.js')
const init = async ()=>{
    const go = new Go()
    const result = await WebAssembly.instantiateStreaming(fetch("/static/nes.wasm"),go.importObject)
    go.run(result.instance)

    console.log("wasm ready for some action yoo")

}

init()

const M_clock = 21_477_272
const CPF = 357_366
const FPS = M_clock/CPF
const MPF = 1000 / FPS

let accum = 0
let last = null 
let running = false

const FBuf = new SharedArrayBuffer(256*240*4)
const FBytes = new Uint8Array(FBuf)

function loop(current){
    if (!running){
        return 
    }

    if (last == null){
        last = current
    }

    elapsed = current - last
    last = current
    accum += elapsed

    accum = Math.min(accum,MPF*4)

    while (accum >= MPF){
        nesFrame()
        accum -= MPF
    }

    const frame = new Uint8Array(frameBuffer.buffer)
    FBytes.set(frame)

    self.postMessage({type:"frameUp"})

    self.requestAnimationFrame(loop)
}

self.onmessage = ({data}) =>{
    switch (data.type){
        case 'init':
            startEmulator()
            console.log()
            self.postMessage({type:"init",FBuf})
            break

        case 'loadRom':
            loadRom(data.rom)
            running = true
            self.requestAnimationFrame(loop)
            break
        
        case 'input':
            update(data.action,data.pressed)
            break


    }
}

const loadRom = async (game) => {
    const response = await fetch( game + ".nes?v=2")
    const buffer = await response.arrayBuffer()

    const uint8view = new Uint8Array(buffer)

    initRom(uint8view)
}


