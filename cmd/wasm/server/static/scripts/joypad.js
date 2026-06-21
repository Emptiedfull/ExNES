

const controlState = {
    'a': false,
    'b': false,
    'select': false,
    'start': false,
    'up': false,
    'down': false,
    'left': false,
    'right': false
}

const ControlMap = {
    "joypad-A":0,
    "joypad-B":1,
    "joypad-select":2,
    "joypad-start":3,
    "dpad-up":4,
    "dpad-down":5,
    "dpad-left":6,
    "dpad-right":7,
}

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

// window.addEventListener('keydown', (e) => {
//     // joypad = document.getElementById("joypad")
//     // joypad.classList.add("active")
//     if (keyMap[e.code] !== undefined  && romRunning) {
//         window.update(ControlMap[keyMap[e.code]],true)

//           btn = document.getElementById(keyMap[e.code])
//           PressBtn(btn)
//     }
// })

// window.addEventListener('keyup', (e) => {
//     if (keyMap[e.code] !== undefined && romRunning) {
        
//         window.update(ControlMap[keyMap[e.code]],false)

//           btn = document.getElementById(keyMap[e.code])
//           ReleaseBtn(btn)
       
//     }
// })


document.addEventListener("DOMContentLoaded",async()=>{
    
    // const dpad_left = document.getElementById("dpad-left") 
    // const dpad_right = document.getElementById("dpad-right")
    // const dpad_up = document.getElementById("dpad-up")
    // const dpad_down = document.getElementById("dpad-down")

    // const startBtn = document.getElementById("joypad-start")
    // const selectBtn = document.getElementById("joypad-select")

    // const ABtn = document.getElementById("joypad-A")
    const BBtn = document.getElementById("joypad-B")

    const toggle = document.getElementById("control-update")

    const buttons = document.querySelectorAll(".joypad-button")
    buttons.forEach(element => {
        element.addEventListener("mousedown",async(e)=>{
            
            await PressBtn(BBtn)
        })
    });
})

const PressBtn = async (button)=>{ 
    button.classList.add("active")

}

const ReleaseBtn = async (button) =>{
    button.classList.remove("active")
}

