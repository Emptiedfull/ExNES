class ConsoleDriver extends AudioWorkletProcessor {
    constructor() {
        super()

        this.samples = null
        this.control = null
        this.SIZE = 0


        this.port.onmessage = ({ data }) => {
            this.SIZE = data.SIZE
            this.samples = new Float32Array(data.audioBufS, 0, this.SIZE)
            this.control = new Int32Array(data.audioBufS, this.SIZE * 4, 3)

        }
    }

    process(inputs, outputs) {
        const output = outputs[0][0]

        if (!this.control) {
            output.fill(0)
            return true
        }


        const wp = Atomics.load(this.control, 0)
        const rp = Atomics.load(this.control, 1)

        const available = wp - rp


        if (available < output.length) {
            console.log("requesting")
            Atomics.store(this.control, 2, 1) //WE WANT MOREEEE
            Atomics.notify(this.control, 2)
            output.fill(0)
            return true
        }

        for (let i = 0; i < output.length; i++) {
            output[i] = this.samples[(rp + i) % this.SIZE]
        }

        Atomics.store(this.control, 1, rp + output.length)

        if (available - output.length < this.SIZE / 2) {
            Atomics.store(this.control, 2, 1)
            Atomics.notify(this.control, 2)
        }

        return true
    }
}

registerProcessor('apu-proc', ConsoleDriver)