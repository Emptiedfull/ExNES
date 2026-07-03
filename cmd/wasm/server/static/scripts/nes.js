import { openCap, closeCap, slotCart, activateJoypad, alignCable, pushbtn, turnKnob, wiggleKnob, knobSettings } from "./graphics.js"
import { initConsole, loadRom, state, updateVolume, UpdateSpeed, PauseGame, ResumeGame } from "./driver.js"
import { wait } from "./joypad.js"
import { createModal } from "./modal.js"
import { activateTip, startRandomTipEngine } from "./tooltips.js"
import { makeMockTiles } from "./rewind.js"



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

const powerBtn = document.getElementById("start")
const pauseBtn = document.getElementById("pause")

let romLoaded = ""
let power = false

const roms = [left, middle, right]

document.addEventListener("DOMContentLoaded", async () => {

    await makeMockTiles()

    startRandomTipEngine()
    await setUpKnobs()
    await setUpButtons()
    await initGames()

    await wait(500)

    updateRom()
});

const BeginKnobs = async () => {
    knobSettings["sound"].setting = 4
    knobSettings["angle"] = 180

    await turnKnob("sound", 180)
    await turnKnob("speed", 180)
}


const setUpKnobs = async () => {
    const soundKnob = document.getElementById("sound")
    const soundIncrements = [0, 45, 90, 135, 180, 225, 270, 315, 360]
    soundKnob.addEventListener("click", async () => {

        if (romLoaded == "" || power == false) {
            await wiggleKnob("sound", 45)
            await turnKnob("sound", 0)
            return
        }

        knobSettings["sound"].setting = knobSettings["sound"].setting + 1

        if (knobSettings["sound"].setting >= soundIncrements.length - 1) {
            knobSettings["sound"].setting = 0
        }

        let targetAngle = soundIncrements[knobSettings["sound"].setting]
        if (targetAngle == 0) {
            updateVolume(0)
        } else {
            updateVolume(targetAngle / 180)
        }
        await turnKnob("sound", targetAngle)
    })

    const speedKnob = document.getElementById("speed")
    const speedIncrements = [45, 90, 135, 180, 225, 270, 315, 360]

    speedKnob.addEventListener("click", async () => {
        if (romLoaded == "" || power == false) {
            await wiggleKnob("speed", 45)
            await turnKnob("speed", 0)
            return
        }

        knobSettings["speed"].setting = knobSettings["speed"].setting + 1

        if (knobSettings["speed"].setting >= speedIncrements.length - 1) {
            knobSettings["speed"].setting = 0
        }

        let targetAngle = speedIncrements[knobSettings["speed"].setting]


        UpdateSpeed(1000 * targetAngle / 180)

        await turnKnob("speed", targetAngle)
    })


}


const initGames = async () => {
    let response = await fetch("./games")

    if (response.ok) {
        let data = await response.json()

        data.forEach(element => {
            GamesArray.push(element)
        });

    }
}



function updateRom() {
    let romSelection = GamesArray.slice(currentIndex - 1, currentIndex + 2)

    if (currentIndex == 0) {

        romSelection = [GamesArray.at(-1), GamesArray[0], GamesArray[1]]
    }

    if (currentIndex == GamesArray.length - 1) {
        romSelection = [GamesArray.at(-2), GamesArray.at(-1), GamesArray[0]]
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

    pauseBtn.addEventListener("click", async () => {
        if (pauseBtn.classList.contains("paused")) {
            pauseBtn.classList.remove("paused")
            pauseBtn.classList.add("unpaused")

              pauseBtn.innerHTML = `   <svg viewBox="0 0 24 24" fill="#5a3a1a" stroke="#5a3a1a" stroke-width="1.5" stroke-linecap="square" xmlns="http://www.w3.org/2000/svg">
                                     <rect x="4" y="3" width="4" height="18"/>
                                     <rect x="16" y="3" width="4" height="18"/></svg>`

            await ResumeGame()
        } else {
            pauseBtn.classList.remove("unpaused")
            pauseBtn.classList.add("paused")

              pauseBtn.innerHTML = ` <svg viewBox="0 0 24 24" fill="#5a3a1a" stroke="#5a3a1a" stroke-width="1.5"
                                    stroke-linecap="square">
                                    <polygon points="6,3 20,12 6,21" />
                                </svg>`

               
            await PauseGame()
        }

    })

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
    nav_back.addEventListener("click", async () => {
        move(-1)
        pushbtn("canvas-back")
    })

    const nav_front = document.getElementById("nav-front")
    nav_front.addEventListener("click", async () => {
        move(1)
        pushbtn("canvas-next")
    })

    middle.addEventListener("click", async () => {


        middle.classList.remove("active")
        await wait(50)
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

        await PauseGame()

        await openCap(cap)
        await wait(200)
        overlay.style.display = "flex"
        overlay.style.transition = "all ease 1"
        overlay.style.opacity = 1
    })


}

async function begin(game) {
    initConsole()
    await loadRom(game)

    await BeginKnobs()


    alignCable(1)

    state.romRunning = true
}



function move(direction) {
    let newIdx = currentIndex + direction

    if (newIdx < 0) {

        currentIndex = GamesArray.length - 1
        updateRom()
        return
    }

    if (newIdx >= GamesArray.length) {

        currentIndex = newIdx - GamesArray.length
        updateRom()
        return
    }

    currentIndex = newIdx
    updateRom()
}








