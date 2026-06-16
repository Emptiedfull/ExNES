const go = new Go()
var loaded = false

WebAssembly.instantiateStreaming(fetch("/static/nes.wasm?v=8"),go.importObject).then((result)=>{
    go.run(result.instance)
    loaded = true
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

const setUpButtons = async ()=>{
    btn = document.getElementById("start")
    btn.addEventListener("click",async()=>{

        startEmulator()
        console.log("console started")
        await loadRom()
        renderLoop()
    })

}

const loadRom = async()=>{
    const response = await fetch("static/mario.nes")
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

