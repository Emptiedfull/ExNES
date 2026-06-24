import { openCap,closeCap,slotCart,activateJoypad,alignCable,pushbtn } from "./graphics.js"
import { initConsole,loadRom,state} from "./driver.js"
import { wait } from "./joypad.js"


let currentIndex = 3

const romArray = [{ "name": "Zelda", "id": "zelda", "img": "zelda" }, { "name": "Zelda", "id": "zelda", "img": "zelda" }, { "name": "donkey kong", "id": "donkey", "img": "donkey" }, { "name": "super mario", "id": "mario", "img": "mario" }, { "name": "Ballon Fight", "id": "balloon", "img": "balloon" }, { "name": "Contra", "id": "contra", "img": "contra" }, { "name": "sb3", "id": "sb3", "img": "mario" }]
const middle = document.getElementById("middle")
const left = document.getElementById("left")
const right = document.getElementById("right")
const slot = document.getElementById("slot")
const overlay = document.getElementById("overlay")
const strip = document.getElementById("strip")
const cap = document.getElementById("cartridge")

const powerLed = document.getElementById("power")


let romLoaded = ""
let power = false

const roms = [left, middle, right]

document.addEventListener("DOMContentLoaded", async () => {

    await setUpButtons()


    updateRom()


    middle.addEventListener("click", async () => {
        middle.classList.remove("active")
        await wait(100)
        let span = middle.querySelector("span")
        span.style.display = "none"
        middle.classList.add("rotated")
        romLoaded = middle.id
        await slotCart()

        if (power) {
            activateJoypad()
            await begin(romLoaded)
        }
    })

    cap.addEventListener("click", async () => {

        await openCap(cap)
        await wait(200)
        overlay.style.display = "flex"
        overlay.style.transition = "all ease 1"
        overlay.style.opacity = 1

        if (power) {
            await begin(romLoaded)
        } else {

        }
    })


});



function updateRom() {

    const romSelection = romArray.slice(currentIndex - 1, currentIndex + 2)

    if (romSelection.length == 3) {
        for (let i = 0; i < 3; i++) {
            const element = romSelection[i];
            const rom = roms[i]

            rom.id = element.id

            const rom_title = rom.querySelector("span")
            const img = rom.querySelector("img")

            img.src = "/rom_images/" + element.img + ".png"
            const spine = rom.querySelector(".rom-spine")
            spine.innerText = element.name
            rom_title.innerText = element.name
        }
    }
}

const setUpButtons = async () => {
    const powerBtn = document.getElementById("start")
    powerBtn.addEventListener("click", async () => {
        power = true

        if (romLoaded != "") {
            begin(romLoaded)
            return
        }

        powerLed.classList.remove("off")

        await openCap(cap)
        await wait(200)
        overlay.style.display = "flex"
        overlay.style.transition = "all ease 1"
        overlay.style.opacity = 1
    })


    const nav_back = document.getElementById("nav-back")
    nav_back.addEventListener("click",async()=>{
        move(-1)
        pushbtn("canvas-back")
    })

    const nav_front = document.getElementById("nav-front")
    nav_front.addEventListener("click",async()=>{
        move(1)
        pushbtn("canvas-next")
    })


}



async function begin(game) {

    await initConsole()
    await loadRom(game)


    alignCable(1)

    state.romRunning = true

}


function move(direction) {
    let newIdx = currentIndex + direction
    if (newIdx >= romArray.length || newIdx < 0) {
        return
    }
    currentIndex = newIdx
    updateRom()
}








