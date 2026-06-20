class AudioProc extends AudioWorkletProcessor {
    constructor(){
        super()

        this.samples = new Float32Array(0)
        this.offset = 0

        this.port.onmessage = (e)=>{
            this.samples = new Float32Array(e.data.buffer)
            this.offset = 0
        }

    }

    process(inputs,outputs){
        const output = outputs[0][0]

        for (let i = 0; i < output.length;i++){
            if (this.offset < this.samples.length) {
                output[i] = this.samples[this.offset++]
            } else {
                output[i] = 0
            }
        }

         if (this.offset >= this.samples.length - 128) {
           console.log("HMM NEED MORE")
        }

        return true
    }

}

registerProcessor('audioproc', AudioProc)