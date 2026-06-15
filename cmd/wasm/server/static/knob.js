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

function drawKnob(canvas, angleDeg) {
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

drawKnob(document.getElementById('sound'), 50);
drawKnob(document.getElementById('brightness'), 120);

