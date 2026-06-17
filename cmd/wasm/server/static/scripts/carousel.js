let currentIndex = 3

const romArray = [{ "name": "Zelda", "id": "zelda", "img": "zelda" },{ "name": "Zelda", "id": "zelda", "img": "zelda" }, { "name": "donkey kong", "id": "donkey", "img": "donkey" }, { "name": "super mario", "id": "mario", "img": "mario" }, { "name": "Ballon Fight", "id": "balloon", "img": "balloon" }, { "name": "Contra", "id": "contra", "img": "contra" }, { "name": "sb3", "id": "sb3", "img": "mario" }]
const middle = document.getElementById("middle")
const left = document.getElementById("left")
const right = document.getElementById("right")
const slot = document.getElementById("slot")
const overlay = document.getElementById("overlay")
const strip = document.getElementById("strip")
const cap = document.getElementById("cartridge")

const roms = [left, middle, right]


document.addEventListener("click", (e) => {
    console.log(e.target)
})

function updateRom() {

    romSelection = romArray.slice(currentIndex - 1, currentIndex + 2)

    if (romSelection.length == 3) {
        for (let i = 0; i < 3; i++) {
            const element = romSelection[i];
            const rom = roms[i]

            text = rom.querySelector("span")
            rom.id = element.id
            img = rom.querySelector("img")
            img.src = "static/rom_images/" + element.img + ".png"
            spine = rom.querySelector(".rom-spine")
            spine.innerText = element.name
            text.innerText = element.name
        }
    }
}

const setUpButtons = async () => {
    btn = document.getElementById("start")
    btn.addEventListener("click", async () => {

        powerLed = document.getElementById("power")
        powerLed.classList.remove("off")
        await openCap()
        await wait(200)
        overlay.style.display = "flex"
        overlay.style.transition= "all ease 1"
        overlay.style.opacity = 1
    })

}


async function begin(game) {
    startEmulator()
    console.log("console started")
    await loadRom(game)
    renderLoop()
}


function move(direction) {
    newIdx = currentIndex + direction
    if (newIdx >= romArray.length || newIdx < 0) {
        return
    }

    currentIndex = newIdx
    updateRom()
}


const drawNav = () => {

    var c = document.getElementById('canvas-back')
    var ctx = c.getContext('2d')
    var PS = 16, COLS = 14, ROWS = 14
    var T = { hi: '#ffc040', md: '#b05000', dk: '#3a1400' }
    var G = [0, 0, 0, 0, 0, 0, 0, 'hi', 'hi', 'hi', 'hi', 'hi', 'hi', 0, 0, 0, 0, 0, 0, 0, 'hi', 'md', 'md', 'md', 'md', 'dk', 0, 0, 0, 0, 0, 0, 0, 'hi', 'md', 'md', 'md', 'md', 'dk', 0, 0, 0, 0, 0, 0, 0, 'hi', 'md', 'md', 'md', 'md', 'dk', 0, 0, 0, 0, 0, 0, 0, 'hi', 'md', 'md', 'md', 'md', 'dk', 0, 0, 0, 0, 0, 0, 0, 'hi', 'md', 'md', 'md', 'md', 'dk', 0, 0, 0, 0, 0, 0, 0, 'hi', 'md', 'md', 'md', 'md', 'md', 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'md', 'md', 'md', 'md', 'md', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'md', 'md', 'md', 'md', 'hi', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'md', 'md', 'md', 'md', 'hi', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'md', 'md', 'md', 'md', 'hi', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'md', 'md', 'md', 'md', 'hi', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'md', 'md', 'md', 'md', 'hi', 0, 0, 0, 0, 0, 0, 0, 0, 0, 'dk', 'dk', 'dk', 'dk', 'dk', 'hi', 0]
    for (var i = 0; i < G.length; i++) {
        if (!G[i]) continue
        var r = Math.floor(i / COLS), col = i % COLS
        ctx.fillStyle = T[G[i]]
        ctx.fillRect(col * PS, r * PS, PS, PS)
    }

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



}


document.addEventListener("DOMContentLoaded", async () => {

    drawNav()
    updateRom()


    middle.addEventListener("click", async () => {
        middle.classList.remove("active")
        await wait(200)
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


function pushbtn(canvasId) {

    const c = document.getElementById(canvasId);

    c.style.filter = 'brightness(2)';
    c.style.transform = 'scale(0.85)';

    setTimeout(() => {
        c.style.filter = 'brightness(1.3)';
        c.style.transform = 'scale(1.1)';
    }, 150);

    setTimeout(() => {
        c.style.filter = '';
        c.style.transform = 'scale(1)';
    }, 280);
}


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


    const slotRect = slot.getBoundingClientRect()
    const stripRect = strip.getBoundingClientRect()

    dy = stripRect.left - slotRect.left



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

    romled = document.getElementById("rom")
    flickerLed(romled)

    await wait(500)



    await closeCap()


    overlay.style.display = "none"



}

const closeCap = async () => {
    cap.style.transition = "all ease 1s"
    cap.classList.remove("open")

    await wait(200)
    await begin(middle.id)


}

const openCap = async () => {
    cap.style.transition = "all ease 1s"
    cap.classList.add("open")

    await await (200)
}



function wait(ms) {
    return new Promise(r => setTimeout(r, ms));
}


const flickerLed = async (led) => {

    for (let i = 0; i < 1; i++) {
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
}
