
var currentPC = 0

document.addEventListener("DOMContentLoaded", async () => {
    await fetchUpdate() 
    await setupButtons()

})

const Base = "http://localhost:8080/"

const setupButtons = async () => {

    btn = document.getElementById("start-debugger")
    btn.addEventListener("click", async () => {
        const response = await fetch(Base + "Debugger/reset")

        if (!response.ok) {
            throw new Error("error starting debugger")
        }

        await fetchUpdate()

    })

    beginbtn = document.getElementById("start-execution")
    beginbtn.addEventListener("click", async () => {
        const response = await fetch(Base + "console/start")

        if (!response.ok) {
            throw new Error("error starting console")
        }

        console.log("starting execution....")
    })


    controlbtn = document.getElementById("control-execution")
    controlbtn.addEventListener("click", async () => {
        console.log(controlbtn.classList.contains("pause"))

        if (controlbtn.classList.contains("pause")) {
            controlbtn.classList.add("unpause")
            controlbtn.classList.remove("pause")
            controlbtn.innerText = "unpause"

            const response = await fetch(Base + "console/pause")

            if (!response.ok) {
                throw new Error("unable to pause console")
            }
        } else {
            controlbtn.classList.add("pause")
            controlbtn.classList.remove("unpause")
            controlbtn.innerText = "pause"

            const response = await fetch(Base + "console/unpause")

            if (!response.ok) {
                throw new Error("unable to unpause console")
            }
        }
    })

    cyclebtn = document.getElementById("run-cycle")
    cyclebtn.addEventListener("click", async () => {
        console.log("running cycle")
        const response = await fetch(Base + "run/cycle")
        if (!response.ok) {
            throw new Error("erorr running cycle")
        }
        fetchUpdate()
    })

    framebtn = document.getElementById("run-frame")
    framebtn.addEventListener("click", async () => {
        const response = await fetch(Base + "run/frame")
        console.log("running frame")
        if (!response.ok) {
            throw new Error("error running frame")
        }
        fetchUpdate()
    })

    run30btn = document.getElementById("run-30")
    run30btn.addEventListener("click", async () => {
        const response = await fetch(Base + "run/30frame")
    })
}

const updatePauseBtn = async () => {
    controlbtn = document.getElementById("control-execution")
    try {
        const response = await fetch(Base + "console/getExecStatus")

        if (!response.ok) {
            throw new Error("unable to get console pause status")
        }

        data = await response.json()
        if (data) {
            controlbtn.classList.add("unpause")
            if (controlbtn.classList.contains("pause")) {
                controlbtn.classList.remove("pause")
            }
            controlbtn.innerText = "Unpause"
        } else {
            controlbtn.classList.add("pause")
            if (controlbtn.classList.contains("unpause")) {
                controlbtn.classList.remove("unpause")
            }
            controlbtn.innerText = "Pause"
        }

        controlbtn.disabled = false
    } catch (error) {
        console.log(error)
        controlbtn.disabled = true

    }
}

const getConsoleStatus = async () => {
    try {
        const response = await fetch(Base + "Debugger/getConsoleStatus")
        data = await response.json()
        console.log(data)
        return data
    } catch (error) {
        console.log(error)
        return false
    }
}



const getScreenArr = async () => {
    try {
        const response = await fetch(Base + "screen/get/Debug")

        const buffer = await response.arrayBuffer()

        const Arr = new Uint8Array(buffer)

        return Arr
    } catch (error) {
        console.log(error)
    }
}

const getState = async () => {
    try {

        const response = await fetch("http://localhost:8080/cpu/state")

        if (!response.ok) {
            throw new Error("HHTP ERROR", response)
        }

        const data = await response.json()


        document.getElementById("A").value = data.a
        document.getElementById("X").value = data.x
        document.getElementById("Y").value = data.y
        document.getElementById("PC").value = formatHex(data.pc)
        document.getElementById("Cycles").value = data.cycles

        currentPC = formatHex(data.pc)

        togglecheckbox("Carry", data.flags.carry)
        togglecheckbox("Zero", data.flags.zero)
        togglecheckbox("Interrupt", data.flags.interrupt)
        togglecheckbox("Overflow", data.flags.overflow)
        togglecheckbox("Negative", data.flags.negative)
        togglecheckbox("Decimal", data.flags.decimal)


    } catch (error) {

    }
}

const togglecheckbox = (checkboxID, state) => {

    checkbox = document.getElementById(checkboxID)
    if (state) {
        checkbox.checked = true
    } else {
        checkbox.checked = false
    }
}





const getLines = async () => {
    try {
        const response = await fetch("http://localhost:8080/disassembly")


        if (!response.ok) {
            throw new Error('HTTP ERROR', response)
        }

        const data = await response.json()



        return data

    } catch (error) {
        console.error(error)
    }
}

const populateLines = (Lines) => {
    if (Lines == null) {
        return
    }
    win = document.querySelector(".code-window")
    win.innerHTML = ""

    for (const [addr, line] of Object.entries(Lines)) {

        const lineHTML = `
        <div class="asm-line ${formatHex(addr) == currentPC ? "current-pc" : ""}">
            <span class="addr">${formatHex(addr)}</span>
            <span class="op">${line.disassembly}</span>${Object.hasOwn(line, "val") ? `= <span class="val">${formatHex(line.val)}</span>` : ''}
        </div> `



        win.insertAdjacentHTML('beforeend', lineHTML);
    }

    current = document.querySelector(".current-pc")
    if (current) {
        current.scrollIntoView({
            behavior: "instant",
            block: "center"
        })
    }
    


}

const formatHex = (int) => {
    return parseInt(int).toString(16).toUpperCase()
}

const fetchUpdate = async () => {
    if (await getConsoleStatus()) {
        lines = await getLines()
        await getState()
        await populateLines(lines)
      
        await updateButtonStatus()


    } else {
        console.log("console not initilizaed")
    }

}


const updateButtonStatus = async ()=>{
    await updatePauseBtn()

    beginbtn = document.getElementById("start-execution")
    beginbtn.disabled = false

}