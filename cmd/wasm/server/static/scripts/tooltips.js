import { wait } from "./joypad"
import { createModal } from "./modal"


export const activateTip = (tip) => {
    if (TipsState[tip]["state"]) {
        TipsState[tip]["function"]()
        TipsState[tip]["state"] = false
    }
}

export const startRandomTipEngine = async () => {
    await wait(1000)
    scheduleTip()
}

const scheduleTip = () => {
    let delay = Math.random() * (60000 - 10000) + 10000

    setTimeout(() => {
        let tip = tips[Math.floor(Math.random() * tips.length)];
        TipsState[tip]["function"]()
        scheduleTip()
    }, delay)
}

const tip_fullscreen = () => {
    createModal("Go big or go home", "click on the tv screen to go fullscreen mode (resolution not garunteed)")
}

const tip_knobs = () => {
    createModal("Play with the dials!!", "Adjust the dials on the television to control the sound and speed levels")
}

const tip_swap = () => {
    createModal("You dont need to reload(prolly)", "Click the cartridge slot to change games while the console is running")
}

const tip_feedback = () => {
    createModal("Have suggestions or found errors?", "Dm emptiedfull on slack please")
}

const tip_controls = ()=>{
    createModal("Change controls","Press the update controls button on the joypad")
}



const TipsState = {
    "fullscreen": { "state": true, "function": tip_fullscreen },
    "knobs": { "state": true, "function": tip_knobs },
    "hotswap": { "state": true, "function": tip_swap },
    "feedback": { "state": true, "function": tip_feedback },
    "controls": {}
}

const tips = Object.keys(TipsState)