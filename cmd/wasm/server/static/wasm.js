const go = new Go()
var loaded = false



WebAssembly.instantiateStreaming(fetch("/static/nes.wasm?v=10"),go.importObject).then(async(result)=>{
    go.run(result.instance)
    loaded = true

     window.startEmulator()
      await loadRom('mario')
    renderLoop()
    window.initBuffer(256)

     await setUpAudio()
})


const canvas = document.getElementById("screen")
const ctx = canvas.getContext("2d")
const imageData = ctx.createImageData(256,240)


window.addEventListener("DOMContentLoaded",async ()=>{
    await setUpButtons()
   

    
   
   

    
})

function renderLoop() {
    nesFrame()

    
    imageData.data.set(frameBuffer)
    ctx.putImageData(imageData,0,0)

    requestAnimationFrame(renderLoop)
}

let BufferSize = 256
let buf = new Uint8Array(BufferSize*4)

document.addEventListener('keydown', () => {
    if (audioCtx && audioCtx.state === 'suspended') {
        audioCtx.resume()
    }
}, { once: true })


const setUpAudio = async ()=>{
    console.log("starting audio processor")

    audioCtx = new AudioContext({sampleRate: 44100})
    const proc = audioCtx.createScriptProcessor(BufferSize,0,1)
   
    proc.onaudioprocess = (e)=>{
        window.getSamples(BufferSize,buf)
        e.outputBuffer.getChannelData(0).set(new Float32Array(buf.buffer))
    }

    proc.connect(audioCtx.destination)

}

const loadRom = async(game)=>{
    const response = await fetch("static/supported/"+game+".nes")
    console.log(game)
    const buffer = await response.arrayBuffer()

    const uint8view = new Uint8Array(buffer)

    initRom(uint8view)
}

const controlState = {
    'a': false,
    'b': false,
    'select': false,
    'start': false,
    'up': false,
    'down': false,
    'left': false,
    'right': false
}

const keyMap = {
    'KeyZ': 0,
    'KeyX': 1,
    'ShiftLeft': 2,
    'Enter': 3,
    'ArrowUp': 4,
    'ArrowDown': 5,
    'ArrowLeft': 6,
    'ArrowRight': 7
}


window.addEventListener('keydown', (e) => {
    if (keyMap[e.code] !== undefined && !controlState[keyMap[e.code]] && loaded) {
        controlState[keyMap[e.code]] = true 
        
          window.update(keyMap[e.code],true)
       
    }
})

window.addEventListener('keyup', (e) => {
    if (keyMap[e.code] !== undefined && loaded) {
        controlState[keyMap[e.code]] = false
        
          window.update(keyMap[e.code],false)
       
    }
})

