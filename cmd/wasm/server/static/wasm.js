const go = new Go()
var loaded = false



WebAssembly.instantiateStreaming(fetch("/static/nes.wasm?v=1011"),go.importObject).then(async(result)=>{
    go.run(result.instance)
    loaded = true

    await setUpAudio()

})

const canvas = document.getElementById("screen")
const ctx = canvas.getContext("2d")
const imageData = ctx.createImageData(256,240)


window.addEventListener("DOMContentLoaded",async ()=>{
    await setUpButtons()
    initCables()
    // await wait(200)
    // alignCables(cable,port,wire)

    
})

const initCables= ()=> {
     cable1 = document.getElementById("cable-1")
    port1 = document.getElementById("port-1")
    wire1 = document.getElementById("wire-1")

    cable2 = document.getElementById("cable-2")
    port2 = document.getElementById("port-2")
    wire2 = document.getElementById("wire-2")

    startCable(cable1,port1,wire1)
    startCable(cable2,port2,wire2)

}

const alignCables = (cable,port,wire)=>{

   
    portBox = port.getBoundingClientRect()
    

    cable.style.top = portBox.top + "px"
    cable.style.left = portBox.left + "px"

    bodyBox = document.querySelector("body").getBoundingClientRect()

    console.log(portBox.bottom,bodyBox.bottom) 

    dy = bodyBox.bottom - portBox.bottom + 200
    dx = portBox.width / 4


    wire.style.height = dy + "px"
    wire.style.top = portBox.top + portBox.height / 2 + "px"
    wire.style.left = portBox.left + portBox.width/4 + "px"


    wire.style.width = dx + "px"

}

const startCable = (cable,port,wire)=>{
    portBox = port.getBoundingClientRect()
    bodyBox = document.querySelector("body").getBoundingClientRect()

    cable.style.position = "absolute"
    cable.style.width = portBox.width + "px"
    cable.style.height = portBox.height + "px"
    cable.style.left = portBox.left + "px"
    cable.style.top = bodyBox.bottom + "px"

    wire.style.position = "absolute"
     dy = bodyBox.bottom - portBox.bottom + 200
    dx = portBox.width / 4

    wire.style.height = dy + "px"
    wire.style.top = bodyBox.bottom + portBox.height / 2 + "px"
    wire.style.left = portBox.left + portBox.width/4 + "px"

}

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
    const response = await fetch("static/supported/"+game+".nes?v=2")
    const buffer = await response.arrayBuffer()

    const uint8view = new Uint8Array(buffer)

    initRom(uint8view)
}
