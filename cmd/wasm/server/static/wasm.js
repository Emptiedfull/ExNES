const go = new Go()
WebAssembly.instantiateStreaming(fetch("/static/nes.wasm?v=9990"),go.importObject).then((result)=>{
    go.run(result.instance)
})

const canvas = document.getElementById("screen")
const ctx = canvas.getContext("2d")
const imageData = ctx.createImageData(256,240)


window.addEventListener("DOMContentLoaded",async ()=>{
    console.log("loading")
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
    btn.addEventListener("click",()=>{

        console.log(startEmulator())
    })

    framebtn = document.getElementById("frame")
    framebtn.addEventListener("click",()=>{
        console.log("starting rendering")
        renderLoop()
    })

    rombtn = document.getElementById("rom")
    rombtn.addEventListener("click",async ()=>{
        await loadRom()
    })

    getbtn = document.getElementById("get")
    getbtn.addEventListener("click",async()=>{
        console.log(getFrame())
    })
}

const loadRom = async()=>{
    const response = await fetch("static/mario.nes")
    const buffer = await response.arrayBuffer()

    const uint8view = new Uint8Array(buffer)

    initRom(uint8view)
}