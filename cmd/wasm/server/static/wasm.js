const go = new Go()
WebAssembly.instantiateStreaming(fetch("/static/emulator.wasm"),go.importObject).then((result)=>{
    go.run(result.instance)
  
})



window.addEventListener("DOMContentLoaded",async ()=>{
    console.log("loading")
    btn = document.getElementById("start")
    btn.addEventListener("click",()=>{
        console.log("btn clicked")
        console.log(startEmulator())
    })

    framebtn = document.getElementById("frame")
    framebtn.addEventListener("click",()=>{
        console.log(runFrame())
    })

    rombtn = document.getElementById("rom")
    rombtn.addEventListener("click",async ()=>{
        await loadRom()
    })
})

const loadRom = async()=>{
    const response = await fetch("static/mario.nes")
    const buffer = await response.arrayBuffer()

    const uint8view = new Uint8Array(buffer)

    initRom(uint8view)
}