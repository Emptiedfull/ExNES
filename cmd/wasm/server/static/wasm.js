const go = new Go()
WebAssembly.instantiateStreaming(fetch("/static/emulator.wasm"),go.importObject).then((result)=>{
    go.run(result.instance)
  
})



window.addEventListener("DOMContentLoaded",()=>{
    console.log("loading")
    btn = document.getElementById("start")
    btn.addEventListener("click",()=>{
        console.log("btn clicked")
        console.log(startEmulator())
    })

    framebtn = document.getElementById("frame")
    framebtn.addEventListener("click",()=>{
        runFrame()
    })
})