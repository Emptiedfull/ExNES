
let BufferSize = 2096

const sharedBuf = new SharedArrayBuffer(BufferSize*4+8) 
const ringBufs = new Float32Array(sharedBuf,0,BufferSize)
const ringIdx = new Int32Array(sharedBuf,0,BufferSize*4,2)

const S_Size = 2048
const S_Buf = new Float32Array(S_Size)

let proc,ctx

const setUpAudio = async ()=>{

    audioCtx = new AudioContext({sampleRate: 44100})

    initBuffer(new Float32Array(2048))

    await audioCtx.audioWorklet.addModule("/static/scripts/audioProc.js?v=1102")
    proc = new AudioWorkletNode(audioCtx,"audioproc",{
        outputChannelCount: [1]
    })

    
    proc.connect(audioCtx.destination)


    proc.port.postMessage({sharedBuf,BufferSize: BufferSize})

    proc.port.onmessage = ({data})=>{
        if (data === "iwantsomesamplespls"){
            const count = getSamples(S_Buf)
            if (count == 0){
                return 
            }

            const wp = Atomics.load(ringIdx,0)
            const rp = Atomics.load(ringIdx,1)

            const free = BufferSize - (wp - rp)
            const x = Math.min(count,free)

            if (x < count){
                console.log("soemthing is super wrong, dropped samples:",count-n)
            }

            for (let i =0;i < x;i++){
                ringBufs[(wp+i)%BufferSize] = S_Buf[i]
            }

            Atomics.store(ringIdx,0,wp+x)
        }
    }

}
