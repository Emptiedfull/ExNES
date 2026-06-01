document.addEventListener("DOMContentLoaded", async () => {
    console.log("running")

    lines = await getLines()
    populateLines(lines)
})


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