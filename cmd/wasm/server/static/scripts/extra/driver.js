const worker = new Worker("/static/scripts/extra/emuWorker.js")
let fBytes = null


const setUpAudio = async (audioBufS,SIZE)=>{
    // console.log("setting up audio")
    const audioCtx = new AudioContext({sampleRate:44100})

    await audioCtx.audioWorklet.addModule("/static/scripts/extra/driverWorklet.js")

    const node = new AudioWorkletNode(audioCtx,'apu-proc',{
        outputChannelCount : [1],
    })

    node.port.postMessage({audioBufS,SIZE})
    node.connect(audioCtx.destination)
}


worker.onmessage = async ({data})=>{
   
    switch (data.type){
        case 'init':
            fBytes = new Uint8Array(data.FBuf)
            console.log("received audioBufS byteLength:", data.audioBufS.byteLength)
            await setUpAudio(data.audioBufS,data.SIZE)
            break
        case 'frameUp':
            imageData.data.set(fBytes)
            ctx.putImageData(imageData, 0, 0)
            break
    }
}

const loadRom = (game) => {
    worker.postMessage({ type: 'loadRom', rom: game })
}