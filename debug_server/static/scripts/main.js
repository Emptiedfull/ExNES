
const canvas = document.getElementById('game-canvas');
const ctx = canvas.getContext('2d');
const fpsCounter = document.getElementById('fps-counter');


const screenWidth = 256;
const screenHeight = 240;
const imgData = ctx.createImageData(screenWidth, screenHeight);


let lastCalledTime;
let fpstag = 0;

document.addEventListener("DOMContentLoaded",async()=>{
    await fetchAndRenderFrame()
})

async function fetchAndRenderFrame() {
    try {
        
        const response = await fetch('http://localhost:8080/screen');
        
        if (!response.ok) throw new Error('Network stream error');
        
        const buffer = await response.arrayBuffer();
        const rawPixels = new Uint8Array(buffer);

        
        if (rawPixels.length === 245760) {
            
            
            imgData.data.set(rawPixels);

            
            ctx.putImageData(imgData, 0, 0);
        }

    } catch (err) {
        console.error("Frame compilation dropped:", err);
    }

}



