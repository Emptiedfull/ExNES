let wasm_up = false
console.log("bro?")
importScripts('/static/wasm_exec.js')
const init = async ()=>{
    const go = new Go()
    const result = await WebAssembly.instantiateStreaming(fetch("/static/nes.wasm"),go.importObject)
    go.run(result.instance)

    console.log("wasm ready for some action yoo")

}

init()


const FBuf = new SharedArrayBuffer(256*240*4)
const FBytes = new Uint8Array(FBuf)

const SIZE = 8192

const audioBufS = new SharedArrayBuffer(SIZE*4 + 4 + 4 +4 ) //uh the plus 4 are different pointers, write,read and updateflag

const samples = new Float32Array(audioBufS,0,SIZE)
const control = new Int32Array(audioBufS,SIZE*4,3)

const S_size = 2048
const S_buf = new Float32Array(S_size)

console.log("control  WORKER byteOffset:", control.byteOffset)
console.log("control[2] WORKER byteOffset:", control.byteOffset + 2*4)


self.onmessage = async ({data}) =>{
    switch (data.type){
        case 'init':
        
            startEmulator()
            initBuffer(new Uint8Array(S_buf.buffer))
            initInput(new Int32Array(data.inputBuf))
            self.postMessage({type:"init",audioBufS,FBuf,SIZE,S_size})
            
            break

        case 'loadRom':
            await loadRom(data.rom)
            running = true
            pump()
            break
        
        case 'input':
            update(data.action,data.pressed)
            break


    }
}

const pump = ()=>{
    while (true){
      
        Atomics.wait(control,2,0)
       

        const wp = Atomics.load(control,0)
        const rp = Atomics.load(control,1)

        const free = SIZE - (wp - rp)
        const want = Math.min(free,S_size)

        if (want > 0){
            drive(want)

            for (let i = 0; i < want;i++){
                samples[(wp + i) % SIZE] = S_buf[i]
            }

            Atomics.store(control,0,wp+want)
        }


        FBytes.set(new Uint8Array(frameBuffer.buffer))
        self.postMessage({type:"frameUp"})

        Atomics.store(control,2,0)
        Atomics.notify(control,2) //god pls work this is my 5th rewrite
    }
}

const loadRom = async (game) => {
    const response = await fetch( game + ".nes?v=2")
    const buffer = await response.arrayBuffer()

    const uint8view = new Uint8Array(buffer)

    initRom(uint8view)
}





