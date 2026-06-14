let nametableInterval = null;

async function refreshNametableDisplay() {
    const canvas = document.getElementById('nametable-canvas');
    if (!canvas) return;
    
    const ctx = canvas.getContext('2d');
    
    try {
        const response = await fetch('http://localhost:8080/screen/get/Debug');
        if (!response.ok) return;
        
        const buffer = await response.arrayBuffer();
        const colorIndices = new Uint8Array(buffer);

        console.log(colorIndices)

        if (colorIndices.length !== 61440) return;

        const imgData = ctx.createImageData(canvas.width, canvas.height);
        
        // Classic developer diagnostic grayscale schema mapping variables
        const palette = {
            0: [0, 0, 0, 255],       // Transparent base layer / Black background
            1: [90, 90, 90, 255],    // Dark Gray structural shadows
            2: [170, 170, 170, 255], // Light Gray details
            3: [255, 255, 255, 255]  // Sharp White typography
        };

        // Linear array streaming sequence loop
        for (let i = 0; i < colorIndices.length; i++) {
            const indexValue = colorIndices[i];
            const [r, g, b, a] = palette[indexValue];

            const canvasIdx = i * 4;
            imgData.data[canvasIdx]     = r;
            imgData.data[canvasIdx + 1] = g;
            imgData.data[canvasIdx + 2] = b;
            imgData.data[canvasIdx + 3] = a;
        }

        ctx.putImageData(imgData, 0, 0);

    } catch (err) {
        console.error("Nametable diagnostics pipeline break:", err);
    }
}

async function drawFullNametableMap() {
    const canvas = document.getElementById('nametable-canvas');
    const ctx = canvas.getContext('2d');
    
    try {
        const response = await fetch('http://localhost:8080/screen/get/Debug');
        const buffer = await response.arrayBuffer();
        const colorIndices = new Uint8Array(buffer); // Size 245760 (512 * 480)

        console.log(colorIndices)

        if (colorIndices.length !== 245760) return;

        const imgData = ctx.createImageData(canvas.width, canvas.height);
        
        const palette = {
            0: [16, 44, 52, 255],     // Original Donkey Kong dark slate blue tint background
            1: [184, 56, 140, 255],   // Vibrant pink/magenta (for girders)
            2: [60, 188, 252, 255],   // Light retro cyan highlight
            3: [252, 252, 252, 255]   // Pure crisp alphanumeric white
        };

        for (let i = 0; i < colorIndices.length; i++) {
            const [r, g, b, a] = palette[colorIndices[i]];
            const canvasIdx = i * 4;
            imgData.data[canvasIdx]     = r;
            imgData.data[canvasIdx + 1] = g;
            imgData.data[canvasIdx + 2] = b;
            imgData.data[canvasIdx + 3] = a;
        }

        ctx.putImageData(imgData, 0, 0);
    } catch (err) {
        console.error(err);
    }
}

// Wire up event listeners once DOM structure settles
document.getElementById('fetch-nametable-btn').addEventListener('click', drawFullNametableMap);

