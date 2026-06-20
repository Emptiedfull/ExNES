class AudioProc extends AudioWorkletProcessor {
    constructor(){
        super()

        this.samples = null
        this.indices = null
        this.SIZE = 0

        this.pending = false

        this.port.onmessage = ({data})=>{
           const {sharedBuffer, bufferSize} = data
           this.SIZE = bufferSize
           this.samples = new Float32Array(sharedBuffer,0,bufferSize)
           this.indices = new Int32Array(sharedBuffer,bufferSize*4,2)
        }

    }

    process(inputs,outputs){
        const output = outputs[0][0]

        if (!this.indices) {
            output.fill(0)
            return true
        }

        const wp = Atomics.load(this.indices,0)
        const rp = Atomics.load(this.indices,1)

        const av = wp - rp

        if (!this.pending && avail < this.SIZE / 2){
            this.port.postMessage('iwantsomesamplespls')
            this.pending = true
        }

        if ( avail < output.length){
            output.fill(0)
            console.log("WE NEED MOREEEEE")
            return true
        }

        for (let i = 0; i < output.length;i++){
            output[i] = this.samples[(rp+1) % this.SIZE]
        }

        Atomics.store(this.indices,1,rp+output.length)
        this.pending = false

       
        return true
    }

}

registerProcessor('audioproc', AudioProc)