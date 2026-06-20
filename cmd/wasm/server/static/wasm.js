const go = new Go()
var loaded = false



WebAssembly.instantiateStreaming(fetch("/static/nes.wasm?v=10"),go.importObject).then(async(result)=>{
    go.run(result.instance)
    loaded = true

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



// document.addEventListener('keydown', () => {
//     if (audioCtx && audioCtx.state === 'suspended') {
//         audioCtx.resume()
//     }
// }, { once: true })


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

