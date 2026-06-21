const worker = new Worker("static/scripts/extra/emuWorker.js")
let fBytes = null

let romRunning = false

window.addEventListener("DOMContentLoaded", async () => {
    await setUpButtons()
    initCables()
})

window.addEventListener('keydown', (e) => {
    if (keyMap[e.code] !== undefined  && romRunning) {
          worker.postMessage({ type: 'input', action: ControlMap[keyMap[e.code]], pressed: true })
          btn = document.getElementById(keyMap[e.code])
          PressBtn(btn)
    }
})

window.addEventListener('keyup', (e) => {
    if (keyMap[e.code] !== undefined && romRunning) {
        
        worker.postMessage({ type: 'input', action: ControlMap[keyMap[e.code]], pressed: false })

          btn = document.getElementById(keyMap[e.code])
          ReleaseBtn(btn)
       
    }
})



const canvas = document.getElementById("screen")
const ctx = canvas.getContext("2d")
const imageData = ctx.createImageData(256, 240)


worker.onmessage = ({data})=>{
    switch (data.type){
        case 'init':
            
            fBytes = new Uint8Array(data.FBuf)
            
            
            break
        case 'frameUp':
            
            imageData.data.set(fBytes)
            ctx.putImageData(imageData, 0, 0)
            break

    }
}


const loadRom = (game)=>{
    worker.postMessage({type:'loadRom',rom:game})
    romRunning = true
}



const initCables = () => {
    cable1 = document.getElementById("cable-1")
    port1 = document.getElementById("port-1")
    wire1 = document.getElementById("wire-1")

    cable2 = document.getElementById("cable-2")
    port2 = document.getElementById("port-2")
    wire2 = document.getElementById("wire-2")

    startCable(cable1, port1, wire1)
    startCable(cable2, port2, wire2)

}

const alignCables = (cable, port, wire) => {


    portBox = port.getBoundingClientRect()


    cable.style.top = portBox.top + "px"
    cable.style.left = portBox.left + "px"

    bodyBox = document.querySelector("body").getBoundingClientRect()

    console.log(portBox.bottom, bodyBox.bottom)

    dy = bodyBox.bottom - portBox.bottom + 200
    dx = portBox.width / 4


    wire.style.height = dy + "px"
    wire.style.top = portBox.top + portBox.height / 2 + "px"
    wire.style.left = portBox.left + portBox.width / 4 + "px"


    wire.style.width = dx + "px"

}

const startCable = (cable, port, wire) => {
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
    wire.style.left = portBox.left + portBox.width / 4 + "px"

}
