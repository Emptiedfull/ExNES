import { state,UpdatePress,UpdateRelease } from "./driver.js"
import { openControlPanel } from "./graphics.js"
import { addUpdatePanel } from "./graphics.js"

const updateBtn = document.getElementById("control-update")

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

const reverseLookUp = {}

const initReverse = ()=>{
    let keys = Object.keys(keyMap)
    keys.forEach(element=>{
        reverseLookUp[keyMap[element]] = element
    })
}

const getKeyNames = ()=>{
    
    let keys = Object.keys(keyMap)
    let result = []
    keys.forEach(element => {
        result.push(keyMap[element])
    });
    return result
}

document.addEventListener("DOMContentLoaded",()=>{
    openControlPanel()

    addUpdatePanel("up","Enter")

   handleUpdateListeners()
   initReverse()
   
    updateBtn.addEventListener('click',()=>{
        console.log("updating the control panel")
        openControlPanel()

    })
})

const handleUpdateListeners = ()=>{
     const buttons = document.querySelectorAll(".joypad-button")
    buttons.forEach(element => {
        element.addEventListener("mousedown",async(e)=>{
            console.log("clicked",element.id)
            
        })
    });
}

window.addEventListener('keydown', (e) => {
    console.log(e.code)
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


const PressBtn = async (button)=>{ 
    button.classList.add("active")

}

const ReleaseBtn = async (button) =>{
    button.classList.remove("active")
}

export function wait(ms) {
    return new Promise(r => setTimeout(r, ms));
}

const listeners = {
    "arrow up":function(

    ){},
}


