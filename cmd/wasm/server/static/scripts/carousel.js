let currentIndex = 0

const romArray = [{ "name": "mario", "id": "mario", "img": "x" }, { "name": "zelda", "id": "zelda", "img": "x" }, { "name": "donkey", "id": "donkey", "img": "x" }, { "name": "balloon", "id": "balloon", "img": "x" }, { "name": "contra", "id": "contra", "img": "x" }, { "name": "sb3", "id": "sb3", "img": "x" }]

const middle = document.getElementById("middle")
const left = document.getElementById("left")
const right = document.getElementById("right")
const slot = document.getElementById("slot")
const overlay = document.getElementById("overlay")
const strip = document.getElementById("strip")
const cap = document.getElementById("cartridge")

const roms = [left, middle, right]



function move(direction) {
    newIdx = currentIndex + direction
    if (newIdx >= romArray.length || newIdx < 0) {
        return
    }

    currentIndex = newIdx

    romSelection = romArray.slice(currentIndex - 1, currentIndex + 2)

    if (romSelection.length == 3) {
        for (let i = 0; i < 3; i++) {
            const element = romSelection[i];

            console.log(roms[i], element)

        }
    }

}


(function () {
    var c = document.getElementById('canvas-back');
    console.log("working")
    var ctx = c.getContext('2d');
    var PS = 16, COLS = 14, ROWS = 14;
    var T = { hi: '#d8d4b0', md: '#8a8660', dk: '#2a2810' };
    var G = [0, 0, 0, 0, 0, 0, 0, 'hi', 'hi', 'hi', 'hi', 'hi', 'hi', 0, 0, 0, 0, 0, 0, 0, 'hi', 'md', 'md', 'md', 'md', 'dk', 0, 0, 0, 0, 0, 0, 0, 'hi', 'md', 'md', 'md', 'md', 'dk', 0, 0, 0, 0, 0, 0, 0, 'hi', 'md', 'md', 'md', 'md', 'dk', 0, 0, 0, 0, 0, 0, 0, 'hi', 'md', 'md', 'md', 'md', 'dk', 0, 0, 0, 0, 0, 0, 0, 'hi', 'md', 'md', 'md', 'md', 'dk', 0, 0, 0, 0, 0, 0, 0, 'hi', 'md', 'md', 'md', 'md', 'md', 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'md', 'md', 'md', 'md', 'md', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'md', 'md', 'md', 'md', 'hi', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'md', 'md', 'md', 'md', 'hi', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'md', 'md', 'md', 'md', 'hi', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'md', 'md', 'md', 'md', 'hi', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'md', 'md', 'md', 'md', 'hi', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'dk', 'dk', 'dk', 'dk', 'hi', 0];
    for (var i = 0; i < G.length; i++) {
        if (!G[i]) continue;
        var r = Math.floor(i / COLS), col = i % COLS;
        ctx.fillStyle = T[G[i]];
        ctx.fillRect(col * PS, r * PS, PS, PS);
    }
})();


(function () {
    var c = document.getElementById('canvas-next');
    var ctx = c.getContext('2d');
    var PS = 16, COLS = 14, ROWS = 14;
    var T = { hi: '#ffc040', md: '#b05000', dk: '#3a1400' };
    var G = [0, 'hi', 'hi', 'hi', 'hi', 'hi', 'hi', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'md', 'md', 'md', 'md', 'hi', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'md', 'md', 'md', 'md', 'hi', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'md', 'md', 'md', 'md', 'hi', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'md', 'md', 'md', 'md', 'hi', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'md', 'md', 'md', 'md', 'hi', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'md', 'md', 'md', 'md', 'md', 'hi', 0, 0, 0, 0, 0, 0, 0, 0, 'md', 'md', 'md', 'md', 'md', 'dk', 0, 0, 0, 0, 0, 0, 0, 'hi', 'md', 'md', 'md', 'md', 'dk', 0, 0, 0, 0, 0, 0, 0, 'hi', 'md', 'md', 'md', 'md', 'dk', 0, 0, 0, 0, 0, 0, 0, 'hi', 'md', 'md', 'md', 'md', 'dk', 0, 0, 0, 0, 0, 0, 0, 'hi', 'md', 'md', 'md', 'md', 'dk', 0, 0, 0, 0, 0, 0, 0, 'hi', 'md', 'md', 'md', 'md', 'dk', 0, 0, 0, 0, 0, 0, 0, 'hi', 'dk', 'dk', 'dk', 'dk', 'dk', 0, 0, 0, 0, 0, 0, 0];
    for (var i = 0; i < G.length; i++) {
        if (!G[i]) continue;
        var r = Math.floor(i / COLS), col = i % COLS;
        ctx.fillStyle = T[G[i]];
        ctx.fillRect(col * PS, r * PS, PS, PS);
    }
})();


document.addEventListener("DOMContentLoaded", async () => {

    middle.addEventListener("click", () => {
        slotCart()
        if (middle.classList.contains("rotated")) {
          
            span = middle.querySelector("span")
            span.style.display = "block"
            middle.classList.remove("rotated")
        } else {
            
            span = middle.querySelector("span")
            span.style.display = "none"
            middle.classList.add("rotated")
        }
    })
});


const slotCart = async () => {
    await wait(800);
       overlay.style.opacity = 0

    const spine = middle.querySelector(".rom-spine");


    const spineRect = spine.getBoundingClientRect();

    const clone = spine.cloneNode(true);
    const computedStyle = window.getComputedStyle(spine);


    clone.style.font = computedStyle.font;
    clone.style.color = computedStyle.color;
    clone.style.letterSpacing = computedStyle.letterSpacing;
    clone.style.textTransform = computedStyle.textTransform;
    clone.style.position = 'fixed';
    clone.style.top = spineRect.top + 'px';
    clone.style.left = spineRect.left + 'px';
    clone.style.width = spineRect.width + 'px';
    clone.style.height = slot.height + 'px';
    clone.style.transition = 'none';
    clone.style.transform = 'none';
    clone.style.margin = '0';
    clone.style.zIndex = '99';

    document.body.appendChild(clone);


    const slotRect =  slot.getBoundingClientRect()
    const stripRect = strip.getBoundingClientRect()

    dy = stripRect.left -slotRect.left 

   

    clone.style.transition = "all 1s ease"
    clone.style.top = slotRect.top + "px"
    clone.style.left = slotRect.left + "px"
    clone.style.width = dy + "px"

    

    await wait(1000)
    clone.style.transition = 'all 0.25s ease-in';
    clone.style.transform = ' scale(0.88)';
    clone.style.transformOrigin = 'center top'; 

      await wait(500);

   
    clone.style.transition = 'all 0.15s ease-out';
    clone.style.transform = 'scale(0.9)';

    await wait(500)

    await closeCap()


    overlay.style.display="none"



}

const closeCap = async ()=>{
    cap.style.transition = "all ease 1s"
   cap.classList.remove("open")

   await wait(200)

   romled = document.getElementById("rom")
   flickerLed(romled)
}



function wait(ms) {
    return new Promise(r => setTimeout(r, ms));
}

move(1)

const flickerLed = async (led) => {
    
    for (let i = 0; i < 6; i++) {
        led.classList.remove("off");
        await wait(60);
        led.classList.add("off");
        await wait(80);
    }

   
    led.classList.remove("off");
    await wait(150);
    led.classList.add("off");
    await wait(120);

    led.classList.remove("off");
    await wait(80);
    led.classList.add("off");
    await wait(200);

   
    led.classList.remove("off");
}
