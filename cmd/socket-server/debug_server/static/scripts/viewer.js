async function renderCHRTable() {
    const canvas = document.getElementById('chr-canvas'); // Must be width="128" height="256"
    const ctx = canvas.getContext('2d');
    
    try {
        const response = await fetch("http://localhost:8080/screen/get/Debug")
        const buffer = await response.arrayBuffer();
        
        // This is your flat []uint8 array coming from Go (size 32768)
        const debugBuffer = new Uint8Array(buffer); 

        const imgData = ctx.createImageData(canvas.width, canvas.height);
        const tilesPerRow = 16; // 128px canvas width / 8px per tile = 16 columns

        // Iterate through all 512 tiles in the stream
        for (let t = 0; t < 512; t++) {
            const tilePixelOffset = t * 64;

            // CRITICAL: Find where this specific 8x8 tile's TOP-LEFT corner 
            // sits on our master 128x256 canvas grid
            const tileX = (t % tilesPerRow) * 8;
            const tileY = Math.floor(t / tilesPerRow) * 8;

            // Loop through the 8 rows and 8 columns of this single tile
            for (let y = 0; y < 8; y++) {
                for (let x = 0; x < 8; x++) {
                    
                    // 1. Grab the color index from your exact Go struct format
                    const colorIndex = debugBuffer[tilePixelOffset + (y * 8 + x)];
                    
                    // 2. Map to your grayscale palette colors
                    let r = 0, g = 0, b = 0;
                    if (colorIndex === 1) { r = 100; g = 100; b = 100; }
                    else if (colorIndex === 2) { r = 180; g = 180; b = 180; }
                    else if (colorIndex === 3) { r = 255; g = 255; b = 255; }

                    // 3. Calculate the absolute pixel coordinates on the 2D canvas sheet
                    const pixelCanvasX = tileX + x;
                    const pixelCanvasY = tileY + y;

                    // 4. Translate that 2D coordinate into the 1D HTML ImageData index
                    const canvasBufferIdx = (pixelCanvasY * canvas.width + pixelCanvasX) * 4;

                    // Write the RGBA data channels safely
                    imgData.data[canvasBufferIdx]     = r;   // R
                    imgData.data[canvasBufferIdx + 1] = g;   // G
                    imgData.data[canvasBufferIdx + 2] = b;   // B
                    imgData.data[canvasBufferIdx + 3] = 255; // A (Fully visible)
                }
            }
        }

        // Flush the perfectly structured matrix to the UI view
        ctx.putImageData(imgData, 0, 0);

    } catch (error) {
        console.error("Failed to map tile alignment values:", error);
    }
}

document.getElementById('fetch-chr-btn').addEventListener('click', renderCHRTable);