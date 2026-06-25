import { openCap,closeCap,slotCart,activateJoypad,alignCable,pushbtn } from "./graphics.js"
import { initConsole,loadRom,state} from "./driver.js"
import { wait } from "./joypad.js"


let currentIndex = 3

const middle = document.getElementById("middle")
const left = document.getElementById("left")
const right = document.getElementById("right")
const slot = document.getElementById("slot")
const overlay = document.getElementById("overlay")
const strip = document.getElementById("strip")
const cap = document.getElementById("cartridge")

const GamesArray = []

const powerLed = document.getElementById("power")


let romLoaded = ""
let power = false

const roms = [left, middle, right]

document.addEventListener("DOMContentLoaded", async () => {

    await setUpButtons()
    await initGames()


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


const initGames = async ()=>{
    let response = await fetch("./games")

    if (response.ok){
        let data = await response.json()

        data.forEach(element => {
            GamesArray.push(element)
        });
       
    }

    console.log(GamesArray)

}



function updateRom() {
     let romSelection = GamesArray.slice(currentIndex - 1, currentIndex + 2)

    if (currentIndex == 0){
         romSelection = [GamesArray.at(-1),GamesArray[0],GamesArray[1]]
    }

    if (romSelection.length == 3) {
        for (let i = 0; i < 3; i++) {
            const element = romSelection[i];
            const rom = roms[i]

            rom.id = element.ID

            const rom_title = rom.querySelector("span")
            const img = rom.querySelector("img")

            img.src = "/rom_images/" + element.ID + ".webp"
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

    if (newIdx < 0 ){
        currentIndex = GamesArray.length + newIdx
        console.log("reseted index to:",currentIndex)
        updateRom()
        return
    }

    if (newIdx >= GamesArray.length){
        currentIndex = currentIndex - GamesArray.length
        updateRom()
        return
    }
    currentIndex = newIdx
    updateRom()
}








