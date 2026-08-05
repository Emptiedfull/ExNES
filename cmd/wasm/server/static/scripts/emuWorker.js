

let wasm_up = false
let rom_loaded = false
let instance = null
importScripts('/wasm_exec.js')
const init = async ()=>{
    const go = new Go()
    const result = await WebAssembly.instantiateStreaming(fetch("/nes.wasm"),go.importObject)
    go.run(result.instance)

    instance = result.instance

    self.postMessage({"type":"wasm"})

      reset()
}

init()

const FBuf = new SharedArrayBuffer(256*240*4)
const FBytes = new Uint8Array(FBuf)

const SIZE = 8192

const audioBufS = new SharedArrayBuffer(SIZE*4 + 4 + 4 +4 + 4 ) //uh the plus 4 are different pointers, write,read and updateflag

const samples = new Float32Array(audioBufS,0,SIZE)
const AudioControl = new Int32Array(audioBufS,SIZE*4,3) // 0,1 buf || 

const gameBuf = new SharedArrayBuffer(4)
const GameControl = new Int32Array(gameBuf) // 2-end Loop, 3-reset 

const frameSigBuf = new SharedArrayBuffer(12)
const frameSig = new Int32Array(frameSigBuf)
       

const S_size = 1024
const S_buf = new Float32Array(S_size)

self.onmessage = async ({data}) =>{
    
    switch (data.type){
        case 'init':

            startEmulator()
            initBuffer(new Uint8Array(S_buf.buffer))
            initInput(new Int32Array(data.inputBuf))
            initSpeed(new Uint8Array(data.speedBuf))
            self.postMessage({type:"init",audioBufS,FBuf,SIZE,S_size,frameSigBuf,gameBuf})
            break

        case 'loadRom':
             console.log("loading rom")
            await loadRom(data.rom)

            self.postMessage({type:"start"})
           
            break
        case 'pump':
            console.log("pump called")
            pump()
            break
        case 'startRaF':
            
            rafPump()
            break

        
        case 'reset':
            
           
            let x = reset()
            if (x == 1){
                  self.postMessage({type:"start"})
            }
          
            break

        case 'getsnap':
            console.log("getting snapshot list")
            let snapshots = await requestSnapshotList()
            
            self.postMessage({type:"snaps",snaps: snapshots})
            break

        case 'loadSnapshot':
            
            loadSnapshot(data.index)
            break
        
        case 'input':
            update(data.action,data.pressed)
            break


    }
}
const frameSize = 256 * 240 * 4

const requestSnapshotList = async ()=>{
    const headerPTR = getSnapshotList()
    const mem = instance.exports.memory.buffer
    const header = new Uint32Array(mem, headerPTR, 3)

    const snapshotPTR = header[0]
    const len = header[1]
    const imagePTR = header[2]

    const snapshotBYTES = new Uint8Array(instance.exports.memory.buffer, snapshotPTR,len)
    console.log(new TextDecoder().decode(snapshotBYTES))
    const snapshots = JSON.parse(new TextDecoder().decode(snapshotBYTES))

 

   for (let i = 0; i < snapshots.length; i++){
        
        const offset = imagePTR  + i * frameSize
        const raw = new Uint8ClampedArray(mem,offset,frameSize)
        const data = new ImageData(new Uint8ClampedArray(raw),256,240)
        snapshots[i].image = await createImageBitmap(data)

   }


    return snapshots
}

const prepareSnapshotImages = async (bufferPTR,frames)=>{
    
 
    const mem = instance.exports.memory.buffer

    const images = []

    for ( let i = 0; i < frames;i++){
        const offset = bufferPTR + i * frameSize

        const rawBuf = new Uint8ClampedArray(mem,offset,frameSize)

        const data = new ImageData(new Uint8ClampedArray(rawBuf),256,240)
        images.push(createImageBitmap(data))
    }

    return Promise.all(images)

}

const rafPump = ()=>{
    console.log("starting loop")
    while (true){
        
         const done = Atomics.load(frameSig,1)
        Atomics.wait(frameSig,0,done)

        if (Atomics.load(frameSig,2) === 1){  
            Atomics.store(frameSig,2,0)
            break
        }

        runFrame()

        FBytes.set(new Uint8Array(frameBuffer.buffer))
        self.postMessage({type:"frameUp"})

         Atomics.add(frameSig, 1, 1)  
    }
}

const pump = ()=>{
    let block = false
   
    while (!block){
        Atomics.wait(AudioControl,2,0)

        const operation = Atomics.exchange(GameControl,0,0)
        

        switch (operation){
            
            
            case 2:
                console.log("ending loop")
                block = true
                
        }

       

        const wp = Atomics.load(AudioControl,0)
        const rp = Atomics.load(AudioControl,1)

        const free = SIZE - (wp - rp)
        const want = Math.min(free,S_size)


        if (want > 0){
            console.log("driving for the want of:",want)
            drive(want)

            for (let i = 0; i < want;i++){
                samples[(wp + i) % SIZE] = S_buf[i]
            }

            Atomics.store(AudioControl,0,wp+want)
        }


        FBytes.set(new Uint8Array(frameBuffer.buffer))
        self.postMessage({type:"frameUp"})

        Atomics.store(AudioControl,2,0)
        Atomics.notify(AudioControl,2) //god pls work this is my 5th rewrite
    }

    console.log("loop ended")

   
}


const loadRom = async (game) => {
    const response = await fetch("/games/" + game + ".nes")

    if (!response.ok){
        createModal("Unable to fetch Rom",response,true)
        return 
    }

    const buffer = await response.arrayBuffer()

    const uint8view = new Uint8Array(buffer)

    initRom(uint8view)
}


