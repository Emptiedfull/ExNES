document.addEventListener("DOMContentLoaded", async () => {
    console.log("running")

    lines = await getLines()
    state = await getState()
    populateLines(lines)

    screen = await getScreenArr()
    console.log(screen)
})

const Base = "http://localhost:8080/"

const setupButtons = ()=>{
    document.getElementById("start-debugger").addEventListener("click",async ()=>{
        const response = await fetch( Base+ "Debugger/reset")
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

        document.getElementById("A").innerText = data.a
        document.getElementById("X").innerText = data.x
        document.getElementById("Y").innerText = data.y
        document.getElementById("PC").innerText = formatHex(data.pc)

        console.log(data)
        
    }catch(error){
        console.log(error)
    }
}


const getLines = async () => {
    try {
        const response = await fetch("http://localhost:8080/disassembly")

        if (!response.ok) {
            throw new Error('HTTP ERROR', response)
        }

        const data = await response.json()
        console.log(data)

        return data

    } catch (error) {
        console.error(error)
    }
}

const populateLines = (Lines) => {
    win = document.querySelector(".code-window")
    console.log(win)
    for (const [addr, line] of Object.entries(Lines)) {
        console.log(` ${formatHex(addr)} ${line.disassembly}`);
        const lineHTML = `
        <div class="asm-line">
            <span class="addr">${formatHex(addr)}</span>
            <span class="op">${line.disassembly}</span>${Object.hasOwn(line,"val") ? `= <span class="val">${formatHex(line.val)}</span>` : ''}
        </div> `

        win.insertAdjacentHTML('beforeend', lineHTML);
    }
}

const formatHex = (int) => {
    return parseInt(int).toString(16).toUpperCase()
}