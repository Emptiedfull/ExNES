document.addEventListener("DOMContentLoaded", async () => {
     console.log("connecting socket")
    connectSocket()
    
    requestAnimationFrame(renderLoop)
})

const startConsole = async () => {
    const response = await fetch("http://localhost:8080/Debugger/reset")
    if (!response.ok) {
        throw new Error("unable to start debugger")
    }
}

let nextFrameBuffer = null;
let paintedFrames = 0;
let lastFPSCheck = performance.now();

const connectSocket = () => {
    const wsUri = "ws://localhost:8080/screen/socket"
    const websocket = new WebSocket(wsUri)
    websocket.addEventListener("open", () => {
        console.log("socket connected")
    })

    websocket.addEventListener("message", async (e) => {
        
        nextFrameBuffer = await e.data.bytes()
    })
}

const canvas = document.getElementById('game-canvas');
const ctx = canvas.getContext('2d');

const screenWidth = 256;
const screenHeight = 240;
const imgData = ctx.createImageData(screenWidth, screenHeight);


function renderLoop() {
   
    if (nextFrameBuffer !== null) {
        RenderFrame(nextFrameBuffer)
        nextFrameBuffer = null; 
        paintedFrames++;
    }

   
    const now = performance.now();
    if (now - lastFPSCheck >= 1000) {
        console.log(`%cFrontend Canvas Render FPS: ${paintedFrames}`, "color: #00ff00; font-weight: bold;");
      
        if (paintedFrames > 35) {
            console.warn("⚠️ Canvas Throttling Warning: Frontend is dropping or rendering frames too quickly!");
        }
        
        paintedFrames = 0;
        lastFPSCheck = now;
    }

    requestAnimationFrame(renderLoop);
}

function RenderFrame(rawPixels) {
    imgData.data.set(rawPixels);
    ctx.putImageData(imgData, 0, 0);
}

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

const keyMap = {
    'KeyZ': 'a',
    'KeyX': 'b',
    'ShiftLeft': 'select',
    'Enter': 'start',
    'ArrowUp': 'up',
    'ArrowDown': 'down',
    'ArrowLeft': 'left',
    'ArrowRight': 'right'
}


window.addEventListener('keydown', (e) => {
    if (keyMap[e.code] !== undefined && !controlState[keyMap[e.code]]) {
        controlState[keyMap[e.code]] = true 
        console.log("Pressed:", keyMap[e.code]) 
        sendUpdate()
    }
})

window.addEventListener('keyup', (e) => {
    if (keyMap[e.code] !== undefined) {
        controlState[keyMap[e.code]] = false
        console.log("Released:", keyMap[e.code])
        sendUpdate()
    }
})

const sendUpdate = () => {
    fetch('http://localhost:8080/controls/update', {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(controlState)
    }).catch(err => console.log(err))
}