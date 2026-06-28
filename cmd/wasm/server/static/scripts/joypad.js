import { state,UpdatePress,UpdateRelease } from "./driver.js"
import { openControlPanel, updateKey } from "./graphics.js"
import { addUpdatePanel,closeControlPanel } from "./graphics.js"

const updateBtn = document.getElementById("control-update")
const buttons = document.querySelectorAll(".joypad-button")


let UpdatingKey = ""

let updatingControls = false

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

// const reverseLookUp = {}

// const initReverse = ()=>{
//     let keys = Object.keys(keyMap)
//     keys.forEach(element=>{
//         reverseLookUp[keyMap[element]] = element
//     })
// }

const getKeyNames = ()=>{
    let keys = Object.keys(keyMap)
    let result = []
    keys.forEach(element => {
        result.push(keyMap[element])
    });
    return result
}


document.addEventListener("DOMContentLoaded",()=>{
   initReverse()
   
    updateBtn.addEventListener('click',()=>{
        if (updatingControls){
             closeControlPanel()    
             removeUpdateListeners()
        }else{
             openControlPanel()
             handleUpdateListeners()
        }
        updatingControls = !updatingControls
        UpdatingKey = ""
    })
})

let controller = null
const handleUpdateListeners = ()=>{

    controller = new AbortController()
    const {signal} = controller
     
    buttons.forEach(element => {
        element.addEventListener("mousedown",async(e)=>{
            console.log("clicked",element.id)
            addUpdatePanel(element.id,getKeyFromAction(element.id))
            UpdatingKey = element.id
        },{signal})
    });
}

const getKeyFromAction = (action)=>{
    let keys = Object.keys(keyMap)
    keys.forEach((key)=>{
        if (keyMap[key] == action){
            return key
        }
    })
}

const removeUpdateListeners = ()=>{
    if (controller){
         controller.abort()
    }
}

const updateBinding = (action,key) =>{
    exsiting = getKeyFromAction(action)
    keyMap[exsiting] = ""
}


window.addEventListener('keydown', (e) => {

    if (UpdatingKey !== "" && updatingControls){


        updateKey(e.code)
        updateBinding(UpdatingKey,e.code)
    }

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


