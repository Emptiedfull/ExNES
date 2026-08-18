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

const neededControls = ['joypad-A','joypad-B','joypad-select','joypad-start','dpad-up','dpad-down','dpad-left','dpad-right']

const getKeyNames = ()=>{
    let keys = Object.keys(keyMap)
    let result = []
    keys.forEach(element => {
        result.push(keyMap[element])
    });
    return result
}


document.addEventListener("DOMContentLoaded",()=>{
   
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
            console.log(element.id,getKeyFromAction(element.id))
            let x = getKeyFromAction(element.id)
            console.log(x)
            addUpdatePanel(element.id,getKeyFromAction(element.id))
            UpdatingKey = element.id
        },{signal})
    });
}

const getKeyFromAction = (action)=>{
    let keys = Object.keys(keyMap)
    let res = ""
    keys.forEach((key)=>{
        if (keyMap[key] == action){
            res = key
        }
    })
    return res
}

const removeUpdateListeners = ()=>{
    if (controller){
        controller.abort()
    }
}

const updateBinding = (action,newKey) =>{
    let old = structuredClone(keyMap)
    let Keys = Object.keys(keyMap)

    Keys.forEach((a)=>{
        if (keyMap[a] == action){
            delete keyMap[a]
        }
    })

    keyMap[newKey] = action

    checkForMissingKeys()
}

const checkForMissingKeys = ()=>{
    let keys = Object.keys(keyMap)
    let needed = structuredClone(neededControls)
    let maped = []

    keys.forEach((key)=>{
        maped.push(keyMap[key])
    })

    needed.forEach((action)=>{
         let x = document.getElementById(action)
        if (maped.includes(action)){
            if (x.classList.contains("key-missing")){
                x.classList.remove("key-missing")
            }
        }else{
            x.classList.add("key-missing")
        }
    })
}


const checkJoypadBindings =()=> {
   let keys = Object.keys(keyMap)
   let present = []

    keys.forEach((key)=>{
        present.push(keyMap[key])
    })

    neededControls.forEach(element => {
        console.log(element,present.includes(element))
    });

}


window.addEventListener('keydown', (e) => {

    if (UpdatingKey !== "" && updatingControls){


        updateKey(e.code)
        updateBinding(UpdatingKey,e.code)
    }

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




