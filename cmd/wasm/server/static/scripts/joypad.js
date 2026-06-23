import { state,UpdatePress,UpdateRelease } from "./driver.js"


const keyMap = {
    'KeyZ': 'joypad-A',
    'KeyX': 'joypad-B',
    'ShiftLeft': 'joypad-select',
    'Enter': 'joypad-start',
    'ArrowUp': 'dpad-up',
    'ArrowDown': 'dpad-down',
    'ArrowLeft': 'dpad-left',
    'ArrowRight': 'dpad-right'
}

window.addEventListener('keydown', (e) => {
    if (keyMap[e.code] !== undefined  && state.romRunning) {
         UpdatePress(keyMap[e.code])

          let btn = document.getElementById(keyMap[e.code])
          PressBtn(btn)
    }
})

window.addEventListener('keyup', (e) => {
    if (keyMap[e.code] !== undefined && state.romRunning) {
        
        UpdateRelease(keyMap[e.code])

          let btn = document.getElementById(keyMap[e.code])
          ReleaseBtn(btn)
       
    }
})


document.addEventListener("DOMContentLoaded",async()=>{
    
    const toggle = document.getElementById("control-update")

    const buttons = document.querySelectorAll(".joypad-button")
    buttons.forEach(element => {
        element.addEventListener("mousedown",async(e)=>{
            
            
        })
    });
})

const PressBtn = async (button)=>{ 
    button.classList.add("active")

}

const ReleaseBtn = async (button) =>{
    button.classList.remove("active")
}

export function wait(ms) {
    return new Promise(r => setTimeout(r, ms));
}

