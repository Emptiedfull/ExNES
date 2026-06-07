
var currentPC = 0

document.addEventListener("DOMContentLoaded", async () => {
    console.log("running")
   
   
    fetchUpdate()
    console.log("setting up buttons")
    await setupButtons()

  

})

const Base = "http://localhost:8080/"

const setupButtons = async()=>{

    btn = document.getElementById("start-debugger")
    btn.addEventListener("click",async ()=>{
        const response = await fetch( Base+ "Debugger/reset")

        if (!response.ok){
            throw new Error("error starting debugger")
        }

    })

    cyclebtn = document.getElementById("run-cycle")
    cyclebtn.addEventListener("click",async()=>{
        console.log("running cycle")
        const response = await fetch(Base + "run/cycle")
        if (!response.ok){
            throw new Error("erorr running cycle")
        }
        fetchUpdate()
    })

    framebtn = document.getElementById("run-frame")
    framebtn.addEventListener("click",async()=>{
        const response = await fetch(Base + "run/frame")
        console.log("running frame")
        if (!response.ok){
            throw new Error("error running frame")
        }
        fetchUpdate()
    })

    run30btn = document.getElementById("run-30")
    run30btn.addEventListener("click",async()=>{
        const response = await fetch(Base + "run/30frame")
    })
}



const getScreenArr = async()=>{
    try{
        const response = await fetch( Base+ "screen/get/Debug")

        const buffer = await response.arrayBuffer()

        const Arr = new Uint8Array(buffer)

        return Arr
    }catch(error){
        console.log(error)
    }
}

const getState = async () =>{
    try {

        const response = await fetch("http://localhost:8080/cpu/state")

        if (!response.ok){
            throw new Error("HHTP ERROR",response)
        }

        const data = await response.json()


        document.getElementById("A").value = data.a
        document.getElementById("X").value = data.x
        document.getElementById("Y").value = data.y
        document.getElementById("PC").value = formatHex(data.pc)

        currentPC = formatHex(data.pc)

        togglecheckbox("Carry",data.flags.carry)
        togglecheckbox("Zero",data.flags.zero)
        togglecheckbox("Interrupt",data.flags.interrupt)
        togglecheckbox("Overflow",data.flags.overflow)
        togglecheckbox("Negative",data.flags.negative)
        togglecheckbox("Decimal",data.flags.decimal)

        
    }catch(error){
        console.log(error)
    }
}

const togglecheckbox = (checkboxID,state) =>{

    checkbox = document.getElementById(checkboxID)
    if (state){
        checkbox.checked = true
    }else{
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
    if (Lines == null){
        return
    }
    win = document.querySelector(".code-window")
    win.innerHTML = ""

    for (const [addr, line] of Object.entries(Lines)) {
     
        const lineHTML = `
        <div class="asm-line ${formatHex(addr) == currentPC ? "current-pc" : "" }">
            <span class="addr">${formatHex(addr)}</span>
            <span class="op">${line.disassembly}</span>${Object.hasOwn(line,"val") ? `= <span class="val">${formatHex(line.val)}</span>` : ''}
        </div> `



        win.insertAdjacentHTML('beforeend', lineHTML);
    }

    current = document.querySelector(".current-pc")
    if (current){
          current.scrollIntoView({
            behavior: "instant", 
            block: "center"     
        })
    }
  
    
}

const formatHex = (int) => {
    return parseInt(int).toString(16).toUpperCase()
}

const fetchUpdate = async()=>{
    lines = await getLines()
    await populateLines(lines)
    await getState()
}