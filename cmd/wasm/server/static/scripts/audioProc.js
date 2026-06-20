class AudioProc extends AudioWorkletProcessor {
    constructor() {
        super()

        this.samples = null
        this.indices = null
        this.SIZE = 0

        this.pending = false

        this.port.onmessage = ({ data }) => {
            const { sharedBuf, BufferSize } = data
            this.SIZE = BufferSize
            this.samples = new Float32Array(sharedBuf, 0, BufferSize)
            this.indices = new Int32Array(sharedBuf, BufferSize * 4, 2)

        }

    }

    process(inputs, outputs) {
       
        const output = outputs[0][0]



        if (this.indices == null) {
            output.fill(0)
            return true
        }

     

        const wp = Atomics.load(this.indices, 0)
        const rp = Atomics.load(this.indices, 1)

        if (!this.indices || !this.SIZE) {
            output.fill(0)
            return true
        }

        const av = wp - rp

        if (!this.pending && av < this.SIZE / 2) {
            this.port.postMessage('iwantsomesamplespls')
            this.pending = true
        }

        if (av < output.length) {
            output.fill(0)
            console.log("WE NEED MOREEEEE")
            this.pending = false
            return true
        }

        for (let i = 0; i < output.length; i++) {
            output[i] = this.samples[(rp + i) % this.SIZE]
        }

        Atomics.store(this.indices, 1, rp + output.length)
        this.pending = false

        return true
    }

}

registerProcessor('audioproc', AudioProc)