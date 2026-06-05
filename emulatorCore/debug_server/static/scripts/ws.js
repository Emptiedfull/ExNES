document.addEventListener("DOMContentLoaded", async () => {

    await startConsole()
    console.log("connecting socket")
    connectSocket()


})


const startConsole = async () => {
    const response = await fetch("http://localhost:8080/Debugger/reset")

    if (!response.ok) {
        throw new Error("unable to start debugger")
    }
}

const connectSocket = () => {
    const wsUri = "ws://localhost:8080/screen/socket"

    const websocket = new WebSocket(wsUri)

    websocket.addEventListener("open", () => {
        console.log("socket connected")
    })

    websocket.addEventListener("message", async (e) => {
        
        pixelArr = await e.data.bytes()
        RenderFrame(pixelArr)
    })
}

const canvas = document.getElementById('game-canvas');
const ctx = canvas.getContext('2d');

const screenWidth = 256;
const screenHeight = 240;

const imgData = ctx.createImageData(screenWidth, screenHeight);

async function RenderFrame(rawPixels) {



    imgData.data.set(rawPixels);


    ctx.putImageData(imgData, 0, 0);



}
