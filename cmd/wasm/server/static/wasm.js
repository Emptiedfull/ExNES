const go = new Go()
WebAssembly.instantiateStreaming(fetch("/static/assests/emulator.wasm"),go.importObject).then((result)=>{
    go.run(result.instance)
    console.log("go wasm lloading")
})