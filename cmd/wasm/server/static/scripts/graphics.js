import { createModal } from "./modal.js"
import { activateTip } from "./tooltips.js"




const body = document.querySelector("body")
const middle = document.getElementById("middle")
const romled = document.getElementById("rom")
const cap = document.getElementById("cartridge")


const screen = document.getElementById("screen")

const control_cont = document.getElementById("joypad-cont")
const panel = document.getElementById("update-panel")
const joypad = document.getElementById("joypad")

const updateBtn = document.getElementById("control-update")
const keyDisplay = document.getElementById("key-display")
const keyText = document.getElementById("show-text")

export const knobSettings = { "speed": { "angle": 0, "setting": 0 }, "sound": { "angle": 0, "setting": 0 } }

const scaleFactor = 2

function wait(ms) {
  return new Promise(r => setTimeout(r, ms));
}


window.addEventListener("DOMContentLoaded", async () => {
  initCables()
  drawNav()
  initCanvas()


  drawKnob("sound", 0)
  drawKnob("speed", 0)

})

export const addUpdatePanel = async (action, key) => {
  keyText.innerText = action + ":  "

  await updateKey(key)
 
}

let flashAnim = keyDisplay.animate([
  {opacity:1},
  {opacity:0.2},
  {opacity:1}
],{
  duration:1000,
  iterations:Infinity
})

flashAnim.pause()


export const updateKey = async (key) => {
  if (key == undefined || key == ""){
    

    keyDisplay.style.backgroundImage = "none"
    keyDisplay.innerText = "NA"
    keyDisplay.removeAttribute('style')

    flashAnim.play()
    return
  }
   const res = await fetch(`./Spritesheets/${key}.png`)

  if (res.ok) {
    flashAnim.cancel()
    keyDisplay.innerText = ""
    const blob = await res.blob()
    const qwn = await createImageBitmap(blob)

    let h = qwn.height * scaleFactor
    let wFull = (qwn.width - 2) * scaleFactor
    let w = wFull / 2

    keyDisplay.style.height = h + "px"
    keyDisplay.style.width = w + "px"

    keyDisplay.style.backgroundSize = `${wFull}px ${h}px `

    keyDisplay.style.backgroundImage = `url('./Spritesheets/${key}.png')`

    keyDisplay.animate([
      { backgroundPosition: '0px 0px' },
      { backgroundPosition: `-${wFull}px 0px` },
    ], {
      duration: 700,
      easing: 'steps(2)',
      iterations: Infinity,
    })
  }
}


const initCanvas = () => {

  screen.addEventListener("mouseenter", async () => {
    activateTip("fullscreen")
  })

  screen.addEventListener("click", async () => {
    try {
      await screen.requestFullscreen();
    } catch (e) {
      await createModal("Unable to go fullscreen", e, true)
    }
  })
}



const initCables = () => {

  let cable1 = document.getElementById("cable-1")
  let port1 = document.getElementById("port-1")
  let wire1 = document.getElementById("wire-1")

  let cable2 = document.getElementById("cable-2")
  let port2 = document.getElementById("port-2")
  let wire2 = document.getElementById("wire-2")


  startCable(cable1, port1, wire1)
  startCable(cable2, port2, wire2)

}

export const alignCable = (cableID) => {
  if (cableID == 1) {
    let cable1 = document.getElementById("cable-1")
    let port1 = document.getElementById("port-1")
    let wire1 = document.getElementById("wire-1")

    alignCables(cable1, port1, wire1)
  } else if (cableID == 2) {
    let cable2 = document.getElementById("cable-2")
    let port2 = document.getElementById("port-2")
    let wire2 = document.getElementById("wire-2")

    alignCables(cable2, port2, wire2)
  }
}

const alignCables = (cable, port, wire) => {


  let portBox = port.getBoundingClientRect()


  cable.style.top = portBox.top + "px"
  cable.style.left = portBox.left + "px"

  let bodyBox = body.getBoundingClientRect()

  let dy = bodyBox.bottom - portBox.bottom + 200
  let dx = portBox.width / 4


  wire.style.height = dy + "px"
  wire.style.top = portBox.top + portBox.height / 2 + "px"
  wire.style.left = portBox.left + portBox.width / 4 + "px"


  wire.style.width = dx + "px"

}

const startCable = (cable, port, wire) => {
  let portBox = port.getBoundingClientRect()
  let bodyBox = document.querySelector("body").getBoundingClientRect()

  cable.style.position = "absolute"
  cable.style.width = portBox.width + "px"
  cable.style.height = portBox.height + "px"
  cable.style.left = portBox.left + "px"
  cable.style.top = bodyBox.bottom + "px"

  wire.style.position = "absolute"
  let dy = bodyBox.bottom - portBox.bottom + 200
  let dx = portBox.width / 4

  wire.style.height = dy + "px"
  wire.style.top = bodyBox.bottom + portBox.height / 2 + "px"
  wire.style.left = portBox.left + portBox.width / 4 + "px"

}

export const activateJoypad = () => {

  let cont = document.getElementById("joypad-cont")

  cont.classList.add("active")
}

export const openControlPanel = () => {
  control_cont.classList.add("updating")
  panel.classList.add("active")

  panel.style.pointerEvents = "all"

  updateBtn.innerText = "Save Controls"
}

export const closeControlPanel = ()=>{
  control_cont.classList.remove("updating")
  panel.classList.remove("active")

  panel.style.pointerEvents = "none"

  updateBtn.innerText = "Update Controls"
}


const C = { //this color pallete was ai generated
  body: '#8B6F4E',
  bodyMid: '#7A5F3E',
  hi: '#C4A882',
  hiTop: '#D9C4A0',
  shadow: '#4A3522',
  shadowD: '#332614',
  rim: '#3B2A16',
  rimHi: '#6B5030',
  dot: '#1E1208',
  dotHi: '#3A2810',
};

function drawKnob(id, angleDeg) {


  knobSettings[id].angle = angleDeg


  const canvas = document.getElementById(id)
  const ctx = canvas.getContext('2d');
  const width = canvas.width
  const height = canvas.height

  const cx = Math.floor(width / 2)
  const cy = Math.floor(height / 2)
  const r = Math.floor(width / 2) - 1;

  const fill = (x, y, col) => {
    if (x < 0 || y < 0 || x >= width || y >= height) return

    ctx.fillStyle = col;
    ctx.fillRect(x, y, 1, 1);
  };


  for (let py = -r; py <= r; py++) {
    for (let px = -r; px <= r; px++) {
      const dist = Math.sqrt(px * px + py * py)
      if (dist > r + 0.5) continue

      const x = cx + px
      const y = cy + py

      if (dist > r - 1) {
        if (px < 0 && py < 0) {
          fill(x, y, C.rimHi)

        }
        else {
          fill(x, y, C.rim)
        }
        continue
      }


      let col = C.body

      if (px <= -Math.floor(r * 0.1) && py <= -Math.floor(r * 0.1) && dist < r * 0.75) {
        col = C.hi;
      }
      if (px <= -Math.floor(r * 0.3) && py <= -Math.floor(r * 0.3) && dist < r * 0.45) {
        col = C.hiTop
      }
      if (px >= Math.floor(r * 0.2) && py >= Math.floor(r * 0.2)) {
        col = C.bodyMid
      }
      if (px >= Math.floor(r * 0.45) && py >= Math.floor(r * 0.45)) {
        col = C.shadow
      }
      if (px >= Math.floor(r * 0.65) && py >= Math.floor(r * 0.65)) {
        col = C.shadowD
      }

      fill(x, y, col)
    }
  }


  const rad = (angleDeg - 90) * Math.PI / 180;
  const len = r - 3;
  const ex = Math.round(cx + Math.cos(rad) * len);
  const ey = Math.round(cy + Math.sin(rad) * len);
  const steps = Math.max(Math.abs(ex - cx), Math.abs(ey - cy));
  for (let i = 1; i <= steps; i++) {
    const t = i / steps;
    const px = Math.round(cx + (ex - cx) * t);
    const py = Math.round(cy + (ey - cy) * t);

    const d = Math.sqrt((px - cx) ** 2 + (py - cy) ** 2);
    if (d < r - 1) fill(px, py, i < 2 ? C.dotHi : C.dot);
  }
  fill(cx, cy, C.dot);
}


export async function turnKnob(id, targetAngle, durationMs = 400, steps = 20) {
  const startAngle = knobSettings[id].angle
  const diff = targetAngle - startAngle
  const stepTime = durationMs / steps

  for (let i = 1; i <= steps; i++) {
    const t = i / steps
    const eased = t < 0.5 ? 2 * t * t : -1 + (4 - 2 * t) * t
    const angle = startAngle + diff * eased
    drawKnob(id, angle)
    await wait(stepTime)
  }
  drawKnob(id, targetAngle)
}


export async function spinKnob(id, rounds = 1, durationMs = 600) {
  const startAngle = knobSettings[id].angle
  const targetAngle = startAngle + (360 * rounds)
  await turnKnob(id, targetAngle, durationMs, 40)
  knobSettings[id].angle = knobSettings[id].angle % 360
}




export async function wiggleKnob(id, intensity = 15) {
  const current = knobSettings[id].angle
  await turnKnob(id, current + intensity, 80, 6)
  await turnKnob(id, current - intensity, 100, 8)
  await turnKnob(id, current + intensity * 0.5, 70, 5)
  await turnKnob(id, current, 80, 6)
}

export async function bootKnob(id, targetAngle = 60) {
  knobSettings[id].angle = -150
  drawKnob(id, -150)
  await wait(100)
  await turnKnob(id, targetAngle, 500, 30)
}

export async function resetKnob(id, targetAngle = 0) {
  await turnKnob(id, 360)
}


export const closeCap = async () => {
  cap.style.transition = "all ease 1s"
  cap.classList.remove("open")
  await wait(200)

}

export function pushbtn(canvasId) {

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


export const openCap = async () => {
  cap.style.transition = "all ease 1s"
  cap.classList.add("open")

  await await (200)
}

export const slotCart = async () => {
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
  clone.style.transform = 'none';
  clone.style.margin = '0';
  clone.style.zIndex = '99';

  document.body.appendChild(clone);


  const slotRect = slot.getBoundingClientRect()
  const stripRect = strip.getBoundingClientRect()

  let dy = stripRect.left - slotRect.left

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


  flickerLed(romled)

  await wait(500)
  await closeCap()

  cleapUpOverlay()
}


const flickerLed = async (led) => {

  led.classList.remove("off");
  await wait(60);
  led.classList.add("off");
  await wait(80);



  led.classList.remove("off");
  await wait(150);
  led.classList.add("off");
  await wait(120);


  led.classList.remove("off");
}


const cleapUpOverlay = () => {
  overlay.style.display = "none"
  middle.classList.remove("rotated")
  middle.classList.add("active")
}


export const drawNav = () => {

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
