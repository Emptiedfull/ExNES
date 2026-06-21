const C = { //this color pallete was ai generated
  body:    '#8B6F4E',
  bodyMid: '#7A5F3E',
  hi:      '#C4A882',
  hiTop:   '#D9C4A0',
  shadow:  '#4A3522',
  shadowD: '#332614',
  rim:     '#3B2A16',
  rimHi:   '#6B5030',
  dot:     '#1E1208',
  dotHi:   '#3A2810',
};

const knobSettings = {"speed":{"angle":60,"setting":1},"sound":{"angle":60,"setting":1}}

function drawKnob(id, angleDeg) {


  knobSettings[id].angle = angleDeg


  const canvas = document.getElementById(id)
  const ctx = canvas.getContext('2d');
  const width = canvas.width
  const height = canvas.height

  const cx = Math.floor(width / 2)
  const cy = Math.floor(height/ 2)
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
    
    const d = Math.sqrt((px-cx)**2 + (py-cy)**2);
    if (d < r - 1) fill(px, py, i < 2 ? C.dotHi : C.dot);
  }
  fill(cx, cy, C.dot);
}

drawKnob("sound", 0)
drawKnob("speed", 0)



async function turnKnob(id, targetAngle, durationMs = 400, steps = 20) {
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


async function spinKnob(id, rounds = 1, durationMs = 600) {
    const startAngle = knobSettings[id].angle
    const targetAngle = startAngle + (360 * rounds)
    await turnKnob(id, targetAngle, durationMs, 40)
    knobSettings[id].angle = knobSettings[id].angle % 360
}


async function stuckKnob(id) {
    const current = knobSettings[id].angle

    
    await turnKnob(id, current + 35, 120, 8)
    await turnKnob(id, current + 20, 60, 5)
    // await turnKnob(id, current + 30, 60, 5)
    // await turnKnob(id, current + 18, 50, 4)
    await turnKnob(id, current + 25, 50, 4)
    await turnKnob(id, current, 150, 10)
}

async function wiggleKnob(id, intensity = 15) {
    const current = knobSettings[id].angle
    await turnKnob(id, current + intensity, 80, 6)
    await turnKnob(id, current - intensity, 100, 8)
    await turnKnob(id, current + intensity * 0.5, 70, 5)
    await turnKnob(id, current, 80, 6)
}

async function bootKnob(id, targetAngle = 60) {
    knobSettings[id].angle = -150
    drawKnob(id, -150)
    await wait(100)
    await turnKnob(id, targetAngle, 500, 30)
}

async function resetKnob(id,targetAngle = 0) {
  await turnKnob(id,360)
}

