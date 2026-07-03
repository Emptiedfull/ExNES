// static/scripts/rewind.js
var track = document.getElementById("items-container");
var preview_img = document.getElementById("preview-img");
var rewind_overlay = document.getElementById("rewind-overlay");
var timeline = document.getElementById("rewind-timeline");
var highligh_view = document.getElementById("rewind-highlight");
var back = document.getElementById("rewind-back");
var next = document.getElementById("rewind-next");
var start = document.getElementById("rewind-start");
var end = document.getElementById("rewind-end");
var load = document.getElementById("review-load");
var cart = document.getElementById("review-cart");
var activeID = 0;
var clicker = new IntersectionObserver(
  (entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        clickSound();
        makeActive(entry.target);
        updateHighlight();
      }
    });
  },
  {
    root: track,
    rootMargin: "0px -45% 0px -45%",
    threshold: 0.1
  }
);
var ctx = new AudioContext();
var createTilesFromSnapshots = async (snapshots) => {
  console.log("creating snapshot tiles");
  let lastTile = null;
  rewind_overlay.style.display = "flex";
  setupRewindButtons(snapshots.length);
  setupTimeline(snapshots.length);
  for (let i = 0; i < snapshots.length; i++) {
    let snapshot = snapshots[i];
    let tile = document.createElement("div");
    tile.id = "tile-" + i;
    clicker.observe(tile);
    tile.classList.add("tape-item");
    let img = document.createElement("canvas");
    img.height = snapshot.image.height;
    img.width = snapshot.image.width;
    tile.addEventListener("mousedown", () => {
      makeActive(tile);
    });
    img.getContext("2d").drawImage(snapshot.image, 0, 0);
    lastTile = tile;
    tile.append(img);
    track.append(tile);
  }
  lastTile.scrollIntoView({
    behavior: "smooth",
    inline: "nearest",
    block: "center"
  });
};
var makeMockTiles = async () => {
  let lastTile = null;
  setupRewindButtons(100);
  setupTimeline(100);
  for (let i = 0; i < 100; i++) {
    let tile = document.createElement("div");
    tile.classList.add("tape-item");
    tile.id = "tile-" + i;
    let img = createRandomPixelArt();
    tile.appendChild(img);
    tile.addEventListener("mousedown", () => {
      makeActive(tile);
    });
    clicker.observe(tile);
    lastTile = tile;
    track.append(tile);
  }
  lastTile.scrollIntoView({
    behavior: "smooth",
    inline: "nearest",
    block: "center"
  });
};
var setupRewindButtons = (length) => {
  next.addEventListener("click", () => {
    let target = parseInt(activeID, 10) + 1;
    if (target > length - 1) {
      target = length - 1;
    }
    let targetEl = document.getElementById("tile-" + target);
    makeActive(document.getElementById("tile-" + target));
    scrollReelToValue(target, convertValToPerc().center);
  });
  back.addEventListener("click", () => {
    let target = parseInt(activeID, 10) - 1;
    if (target < 0) {
      target = 0;
    }
    let targetEl = document.getElementById("tile-" + target);
    makeActive(targetEl);
    scrollReelToValue(target, convertValToPerc().center);
  });
  end.addEventListener("click", async () => {
    let target = length - 1;
    let targetEl = document.getElementById("tile-" + target);
    scrollReelToValue(target, false);
    await wait(200);
    makeActive(targetEl);
  });
  start.addEventListener("click", async () => {
    let target = 0;
    let targetEl = document.getElementById("tile-" + target);
    scrollReelToValue(target, false);
    await wait(200);
    makeActive(targetEl);
  });
  load.addEventListener("click", async () => {
    console.log("loading snapshot:", activeID);
    startSaveLoad();
    rewind_overlay.style.display = "none";
  });
};
var startSaveLoad = () => {
  let clone = cart.cloneNode(true);
  clone.classList.add("rewind-cart-clone");
  let orignalRect = cart.getBoundingClientRect();
  let sourceCanvas = cart.querySelector("canvas");
  let width = orignalRect.right - orignalRect.left;
  clone.style.left = orignalRect.left + "px";
  clone.style.top = orignalRect.top + "px";
  clone.style.width = width + "px";
  clone.querySelector("canvas").getContext("2d").drawImage(sourceCanvas, 0, 0);
  console.log(orignalRect.right, orignalRect.bottom);
  document.body.appendChild(clone);
};
var lastactive = null;
var makeActive = (el) => {
  if (lastactive != null) {
    lastactive.classList.remove("tape-active");
  }
  let source = el.querySelector("canvas");
  preview_img.height = 240;
  preview_img.width = 256;
  preview_img.getContext("2d").drawImage(source, 0, 0);
  el.classList.add("tape-active");
  let index = el.id.slice(5);
  timeline.value = index;
  activeID = index;
  lastactive = el;
};
var setupTimeline = (val) => {
  timeline.min = 0;
  timeline.max = val - 1;
};
timeline.addEventListener("input", (e) => {
  scrollReelToValue(timeline.value, convertValToPerc().center);
});
var updateHighlight = () => {
  let res = convertValToPerc();
  highligh_view.style.left = res.perc - 6 + "%";
  return res.center;
};
var scrollReelToValue = (value, center) => {
  let elemtentID = "tile-" + Math.round(value);
  let tile = document.getElementById(elemtentID);
  tile.scrollIntoView({
    behavior: "instant",
    inline: center ? "center" : "nearest",
    block: "center"
  });
};
var convertValToPerc = () => {
  let max = timeline.max;
  let min = timeline.min;
  let val = timeline.value;
  let perc = (val - min) / (max - min) * 100;
  let center = true;
  if (perc < 6) {
    perc = 6;
    center = false;
  }
  if (perc > 94) {
    perc = 94;
    center = false;
  }
  return { perc, center };
};
var clickSound = () => {
  if (ctx == null) return;
  const osc = ctx.createOscillator();
  const gain2 = ctx.createGain();
  osc.connect(gain2);
  gain2.connect(ctx.destination);
  osc.frequency.value = 900;
  gain2.gain.setValueAtTime(0.4, ctx.currentTime);
  gain2.gain.exponentialRampToValueAtTime(1e-3, ctx.currentTime + 0.02);
  osc.start();
  osc.stop(ctx.currentTime + 0.02);
};
function createRandomPixelArt(blockSize = 4) {
  const canvas2 = document.createElement("canvas");
  canvas2.width = 256;
  canvas2.height = 240;
  const ctx3 = canvas2.getContext("2d");
  const palette = [
    "#000000",
    "#FCFCFC",
    "#F8F8F8",
    "#BCBCBC",
    "#7C7C7C",
    "#A40000",
    "#0000FC",
    "#00A800",
    "#F8B800",
    "#00FCFC",
    "#F800F8",
    "#585858"
  ];
  for (let y = 0; y < canvas2.height; y += blockSize) {
    for (let x = 0; x < canvas2.width; x += blockSize) {
      ctx3.fillStyle = palette[Math.floor(Math.random() * palette.length)];
      ctx3.fillRect(x, y, blockSize, blockSize);
    }
  }
  return canvas2;
}

// static/scripts/driver.js
var worker = new Worker(new URL("./emuWorker.js", import.meta.url));
var fBytes = null;
var inputBuf = new SharedArrayBuffer(4);
var inputState = new Int32Array(inputBuf);
var canvas = document.getElementById("screen");
var ctx2 = canvas.getContext("2d");
var imageData = ctx2.createImageData(256, 240);
var speedBuf = new SharedArrayBuffer(4);
var speedNum = new Int32Array(speedBuf);
var audioBufS = null;
var control = null;
Atomics.store(speedNum, 0, 1e3);
var state = {
  romRunning: false
};
var audioCtx = null;
var gain = null;
var startTime = null;
window.addEventListener("keydown", async (e) => {
  if (e.code == "KeyL") {
    startTime = performance.now();
    await getSnapList();
  } else if (e.code == "KeyP") {
    console.log("hello");
    await loadSnap(0);
  }
});
var setUpAudio = async (audioBufS2, SIZE) => {
  if (audioCtx !== null && gain !== null) {
    return;
  }
  audioCtx = new AudioContext({ sampleRate: 44100 });
  await audioCtx.audioWorklet.addModule(new URL("./driverWorklet.js", import.meta.url));
  const node = new AudioWorkletNode(audioCtx, "apu-proc", {
    outputChannelCount: [1]
  });
  gain = audioCtx.createGain();
  node.port.postMessage({ audioBufS: audioBufS2, SIZE });
  node.connect(gain);
  gain.connect(audioCtx.destination);
};
var PauseGame = async () => {
  if (audioCtx == null) {
    return;
  }
  await audioCtx.suspend();
  Atomics.store(control, 2, 2);
  Atomics.notify(control, 2);
};
var ResumeGame = async () => {
  console.log("playing");
  if (audioCtx == null) {
    return;
  }
  await audioCtx.resume();
  Atomics.store(control, 2, 1);
  Atomics.notify(control, 2);
  worker.postMessage({ "type": "pump" });
};
var getSnapList = async () => {
  await PauseGame();
  worker.postMessage({ type: "getsnap" });
};
var loadSnap = async (index) => {
  worker.postMessage({ type: "loadSnapshot", index });
  await wait(100);
  await ResumeGame();
};
var updateVolume = (intensity) => {
  if (gain !== null) {
    gain.gain.setValueAtTime(gain.gain.value, audioCtx.currentTime);
    gain.gain.linearRampToValueAtTime(intensity, audioCtx.currentTime + 1);
  }
};
var initConsole = () => {
  worker.postMessage({ type: "init", speedBuf, inputBuf });
};
worker.onmessage = async ({ data }) => {
  switch (data.type) {
    case "init":
      fBytes = new Uint8Array(data.FBuf);
      audioBufS = data.audioBufS;
      control = new Int32Array(audioBufS, data.SIZE * 4, 3);
      await setUpAudio(data.audioBufS, data.SIZE);
      break;
    case "wasm":
      createModal("Console Ready!!", "Press the power button to start the console");
      worker.postMessage({ type: "init", inputBuf, speedBuf });
      break;
    case "snaps":
      await createTilesFromSnapshots(data.snaps);
      break;
    case "frameUp":
      imageData.data.set(fBytes);
      ctx2.putImageData(imageData, 0, 0);
      break;
  }
};
var ControlMap = {
  "joypad-A": 0,
  "joypad-B": 1,
  "joypad-select": 2,
  "joypad-start": 3,
  "dpad-up": 4,
  "dpad-down": 5,
  "dpad-left": 6,
  "dpad-right": 7
};
var loadRom = (game) => {
  Atomics.store(control, 2, 1);
  Atomics.notify(control, 2);
  audioCtx.resume();
  worker.postMessage({ type: "loadRom", rom: game });
};
var UpdateSpeed = (speed) => {
  Atomics.store(speedNum, 0, speed);
};
var UpdatePress = (btn) => {
  Atomics.or(inputState, 0, 1 << ControlMap[btn]);
};
var UpdateRelease = (btn) => {
  Atomics.and(inputState, 0, ~(1 << ControlMap[btn]));
};

// static/scripts/joypad.js
var updateBtn = document.getElementById("control-update");
var buttons = document.querySelectorAll(".joypad-button");
var UpdatingKey = "";
var updatingControls = false;
var keyMap = {
  "KeyZ": "joypad-A",
  "KeyX": "joypad-B",
  "ShiftLeft": "joypad-select",
  "Enter": "joypad-start",
  "ArrowUp": "dpad-up",
  "ArrowDown": "dpad-down",
  "ArrowLeft": "dpad-left",
  "ArrowRight": "dpad-right"
};
var neededControls = ["joypad-A", "joypad-B", "joypad-select", "joypad-start", "dpad-up", "dpad-down", "dpad-left", "dpad-right"];
document.addEventListener("DOMContentLoaded", () => {
  updateBtn.addEventListener("click", () => {
    if (updatingControls) {
      closeControlPanel();
      removeUpdateListeners();
    } else {
      openControlPanel();
      handleUpdateListeners();
    }
    updatingControls = !updatingControls;
    UpdatingKey = "";
  });
});
var controller = null;
var handleUpdateListeners = () => {
  controller = new AbortController();
  const { signal } = controller;
  buttons.forEach((element) => {
    element.addEventListener("mousedown", async (e) => {
      console.log(element.id, getKeyFromAction(element.id));
      let x = getKeyFromAction(element.id);
      console.log(x);
      addUpdatePanel(element.id, getKeyFromAction(element.id));
      UpdatingKey = element.id;
    }, { signal });
  });
};
var getKeyFromAction = (action) => {
  let keys = Object.keys(keyMap);
  let res = "";
  keys.forEach((key) => {
    if (keyMap[key] == action) {
      res = key;
    }
  });
  return res;
};
var removeUpdateListeners = () => {
  if (controller) {
    controller.abort();
  }
};
var updateBinding = (action, newKey) => {
  let old = structuredClone(keyMap);
  let Keys = Object.keys(keyMap);
  Keys.forEach((a) => {
    if (keyMap[a] == action) {
      delete keyMap[a];
    }
  });
  keyMap[newKey] = action;
  checkForMissingKeys();
};
var checkForMissingKeys = () => {
  let keys = Object.keys(keyMap);
  let needed = structuredClone(neededControls);
  let maped = [];
  keys.forEach((key) => {
    maped.push(keyMap[key]);
  });
  needed.forEach((action) => {
    let x = document.getElementById(action);
    if (maped.includes(action)) {
      if (x.classList.contains("key-missing")) {
        x.classList.remove("key-missing");
      }
    } else {
      x.classList.add("key-missing");
    }
  });
};
window.addEventListener("keydown", (e) => {
  if (UpdatingKey !== "" && updatingControls) {
    updateKey(e.code);
    updateBinding(UpdatingKey, e.code);
  }
  if (keyMap[e.code] !== void 0 && state.romRunning) {
    UpdatePress(keyMap[e.code]);
    let btn = document.getElementById(keyMap[e.code]);
    PressBtn(btn);
  }
});
window.addEventListener("keyup", (e) => {
  if (keyMap[e.code] !== void 0 && state.romRunning) {
    UpdateRelease(keyMap[e.code]);
    let btn = document.getElementById(keyMap[e.code]);
    ReleaseBtn(btn);
  }
});
var PressBtn = async (button) => {
  button.classList.add("active");
};
var ReleaseBtn = async (button) => {
  button.classList.remove("active");
};
function wait(ms) {
  return new Promise((r) => setTimeout(r, ms));
}

// static/scripts/modal.js
var createModal = async (head, body2, error = false, fleeting = true) => {
  let modalList = document.querySelector(".moodle-bar");
  const modal = document.createElement("div");
  modal.classList.add("modal", "modal-corners");
  const title = document.createElement("h1");
  if (error) {
    title.classList.add("error");
  } else {
    title.classList.add("info");
  }
  title.textContent = head;
  modal.appendChild(title);
  const para = document.createElement("p");
  para.classList.add("modal-text");
  para.textContent = body2;
  modal.appendChild(para);
  modalList.appendChild(modal);
  if (fleeting) {
    setUpForFailure(modal);
  }
  activateModal(modal);
  return modal;
};
var activateModal = (modal) => {
  modal.classList.add("active");
};
var deactivateModal = async (modal) => {
  modal.classList.remove("active");
  modal.classList.add("inactive");
  await wait(400);
  modal.remove();
};
var setUpForFailure = async (modal) => {
  await wait(3500);
  if (modal.classList.contains("active")) {
    await deactivateModal(modal);
  } else {
    console.log("bro alr failed");
  }
};

// static/scripts/tooltips.js
var activateTip = (tip) => {
  if (TipsState[tip]["state"]) {
    TipsState[tip]["function"]();
    TipsState[tip]["state"] = false;
  }
};
var startRandomTipEngine = () => {
  scheduleTip();
};
var scheduleTip = () => {
  let delay = Math.random() * (6e4 - 1e4) + 1e4;
  setTimeout(() => {
    let tip = tips[Math.floor(Math.random() * tips.length)];
    TipsState[tip]["function"]();
    scheduleTip();
  }, delay);
};
var tip_fullscreen = () => {
  createModal("Go big or go home", "click on the tv screen to go fullscreen mode (resolution not garunteed)");
};
var tip_knobs = () => {
  createModal("Play with the dials!!", "Adjust the dials on the television to control the sound and speed levels");
};
var tip_swap = () => {
  createModal("You dont need to reload(prolly)", "Click the cartridge slot to change games while the console is running");
};
var tip_feedback = () => {
  createModal("Have suggestions or found errors?", "Dm emptiedfull on slack please");
};
var TipsState = {
  "fullscreen": { "state": true, "function": tip_fullscreen },
  "knobs": { "state": true, "function": tip_knobs },
  "hotswap": { "state": true, "function": tip_swap },
  "feedback": { "state": true, "function": tip_feedback },
  "controls": {}
};
var tips = Object.keys(TipsState);

// static/scripts/graphics.js
var body = document.querySelector("body");
var middle = document.getElementById("middle");
var romled = document.getElementById("rom");
var cap = document.getElementById("cartridge");
var screen = document.getElementById("screen");
var control_cont = document.getElementById("joypad-cont");
var panel = document.getElementById("update-panel");
var joypad = document.getElementById("joypad");
var updateBtn2 = document.getElementById("control-update");
var keyDisplay = document.getElementById("key-display");
var keyText = document.getElementById("show-text");
var knobSettings = { "speed": { "angle": 0, "setting": 0 }, "sound": { "angle": 0, "setting": 0 } };
var scaleFactor = 2;
function wait2(ms) {
  return new Promise((r) => setTimeout(r, ms));
}
window.addEventListener("DOMContentLoaded", async () => {
  initCables();
  drawNav();
  initCanvas();
  drawKnob("sound", 0);
  drawKnob("speed", 0);
});
var addUpdatePanel = async (action, key) => {
  keyText.innerText = action + ":  ";
  await updateKey(key);
};
var flashAnim = keyDisplay.animate([
  { opacity: 1 },
  { opacity: 0.2 },
  { opacity: 1 }
], {
  duration: 1e3,
  iterations: Infinity
});
flashAnim.pause();
var updateKey = async (key) => {
  if (key == void 0 || key == "") {
    keyDisplay.style.backgroundImage = "none";
    keyDisplay.innerText = "NA";
    keyDisplay.removeAttribute("style");
    flashAnim.play();
    return;
  }
  const res = await fetch(`./dist/spritesheets/${key}.png`);
  if (res.ok) {
    flashAnim.cancel();
    keyDisplay.innerText = "";
    const blob = await res.blob();
    const qwn = await createImageBitmap(blob);
    let h = qwn.height * scaleFactor;
    let wFull = (qwn.width - 2) * scaleFactor;
    let w = wFull / 2;
    keyDisplay.style.height = h + "px";
    keyDisplay.style.width = w + "px";
    keyDisplay.style.backgroundSize = `${wFull}px ${h}px `;
    keyDisplay.style.backgroundImage = `url('./dist/spritesheets/${key}.png')`;
    keyDisplay.animate([
      { backgroundPosition: "0px 0px" },
      { backgroundPosition: `-${wFull}px 0px` }
    ], {
      duration: 700,
      easing: "steps(2)",
      iterations: Infinity
    });
  }
};
var initCanvas = () => {
  screen.addEventListener("mouseenter", async () => {
    activateTip("fullscreen");
  });
  screen.addEventListener("click", async () => {
    try {
      await screen.requestFullscreen();
    } catch (e) {
      await createModal("Unable to go fullscreen", e, true);
    }
  });
};
var initCables = () => {
  let cable1 = document.getElementById("cable-1");
  let port1 = document.getElementById("port-1");
  let wire1 = document.getElementById("wire-1");
  let cable2 = document.getElementById("cable-2");
  let port2 = document.getElementById("port-2");
  let wire2 = document.getElementById("wire-2");
  startCable(cable1, port1, wire1);
  startCable(cable2, port2, wire2);
};
var alignCable = (cableID) => {
  if (cableID == 1) {
    let cable1 = document.getElementById("cable-1");
    let port1 = document.getElementById("port-1");
    let wire1 = document.getElementById("wire-1");
    alignCables(cable1, port1, wire1);
  } else if (cableID == 2) {
    let cable2 = document.getElementById("cable-2");
    let port2 = document.getElementById("port-2");
    let wire2 = document.getElementById("wire-2");
    alignCables(cable2, port2, wire2);
  }
};
var alignCables = (cable, port, wire) => {
  let portBox = port.getBoundingClientRect();
  cable.style.top = portBox.top + "px";
  cable.style.left = portBox.left + "px";
  let bodyBox = body.getBoundingClientRect();
  let dy = bodyBox.bottom - portBox.bottom + 200;
  let dx = portBox.width / 4;
  wire.style.height = dy + "px";
  wire.style.top = portBox.top + portBox.height / 2 + "px";
  wire.style.left = portBox.left + portBox.width / 4 + "px";
  wire.style.width = dx + "px";
};
var startCable = (cable, port, wire) => {
  let portBox = port.getBoundingClientRect();
  let bodyBox = document.querySelector("body").getBoundingClientRect();
  cable.style.position = "absolute";
  cable.style.width = portBox.width + "px";
  cable.style.height = portBox.height + "px";
  cable.style.left = portBox.left + "px";
  cable.style.top = bodyBox.bottom + "px";
  wire.style.position = "absolute";
  let dy = bodyBox.bottom - portBox.bottom + 200;
  let dx = portBox.width / 4;
  wire.style.height = dy + "px";
  wire.style.top = bodyBox.bottom + portBox.height / 2 + "px";
  wire.style.left = portBox.left + portBox.width / 4 + "px";
};
var activateJoypad = () => {
  let cont = document.getElementById("joypad-cont");
  cont.classList.add("active");
};
var openControlPanel = () => {
  control_cont.classList.add("updating");
  panel.classList.add("active");
  panel.style.pointerEvents = "all";
  updateBtn2.innerText = "Save Controls";
};
var closeControlPanel = () => {
  control_cont.classList.remove("updating");
  panel.classList.remove("active");
  panel.style.pointerEvents = "none";
  updateBtn2.innerText = "Update Controls";
};
var C = {
  //this color pallete was ai generated
  body: "#8B6F4E",
  bodyMid: "#7A5F3E",
  hi: "#C4A882",
  hiTop: "#D9C4A0",
  shadow: "#4A3522",
  shadowD: "#332614",
  rim: "#3B2A16",
  rimHi: "#6B5030",
  dot: "#1E1208",
  dotHi: "#3A2810"
};
function drawKnob(id, angleDeg) {
  knobSettings[id].angle = angleDeg;
  const canvas2 = document.getElementById(id);
  const ctx3 = canvas2.getContext("2d");
  const width = canvas2.width;
  const height = canvas2.height;
  const cx = Math.floor(width / 2);
  const cy = Math.floor(height / 2);
  const r = Math.floor(width / 2) - 1;
  const fill = (x, y, col) => {
    if (x < 0 || y < 0 || x >= width || y >= height) return;
    ctx3.fillStyle = col;
    ctx3.fillRect(x, y, 1, 1);
  };
  for (let py = -r; py <= r; py++) {
    for (let px = -r; px <= r; px++) {
      const dist = Math.sqrt(px * px + py * py);
      if (dist > r + 0.5) continue;
      const x = cx + px;
      const y = cy + py;
      if (dist > r - 1) {
        if (px < 0 && py < 0) {
          fill(x, y, C.rimHi);
        } else {
          fill(x, y, C.rim);
        }
        continue;
      }
      let col = C.body;
      if (px <= -Math.floor(r * 0.1) && py <= -Math.floor(r * 0.1) && dist < r * 0.75) {
        col = C.hi;
      }
      if (px <= -Math.floor(r * 0.3) && py <= -Math.floor(r * 0.3) && dist < r * 0.45) {
        col = C.hiTop;
      }
      if (px >= Math.floor(r * 0.2) && py >= Math.floor(r * 0.2)) {
        col = C.bodyMid;
      }
      if (px >= Math.floor(r * 0.45) && py >= Math.floor(r * 0.45)) {
        col = C.shadow;
      }
      if (px >= Math.floor(r * 0.65) && py >= Math.floor(r * 0.65)) {
        col = C.shadowD;
      }
      fill(x, y, col);
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
async function turnKnob(id, targetAngle, durationMs = 400, steps = 20) {
  const startAngle = knobSettings[id].angle;
  const diff = targetAngle - startAngle;
  const stepTime = durationMs / steps;
  for (let i = 1; i <= steps; i++) {
    const t = i / steps;
    const eased = t < 0.5 ? 2 * t * t : -1 + (4 - 2 * t) * t;
    const angle = startAngle + diff * eased;
    drawKnob(id, angle);
    await wait2(stepTime);
  }
  drawKnob(id, targetAngle);
}
async function wiggleKnob(id, intensity = 15) {
  const current = knobSettings[id].angle;
  await turnKnob(id, current + intensity, 80, 6);
  await turnKnob(id, current - intensity, 100, 8);
  await turnKnob(id, current + intensity * 0.5, 70, 5);
  await turnKnob(id, current, 80, 6);
}
var closeCap = async () => {
  cap.style.transition = "all ease 1s";
  cap.classList.remove("open");
  await wait2(200);
};
function pushbtn(canvasId) {
  const c = document.getElementById(canvasId);
  c.style.filter = "brightness(2)";
  c.style.transform = "scale(0.85)";
  setTimeout(() => {
    c.style.filter = "brightness(1.3)";
    c.style.transform = "scale(1.1)";
  }, 150);
  setTimeout(() => {
    c.style.filter = "";
    c.style.transform = "scale(1)";
  }, 280);
}
var openCap = async () => {
  cap.style.transition = "all ease 1s";
  cap.classList.add("open");
  await await 200;
};
var slotCart = async () => {
  await wait2(800);
  overlay.style.opacity = 0;
  const spine = middle.querySelector(".rom-spine");
  const spineRect = spine.getBoundingClientRect();
  const clone = spine.cloneNode(true);
  const computedStyle = window.getComputedStyle(spine);
  clone.style.font = computedStyle.font;
  clone.style.color = computedStyle.color;
  clone.style.letterSpacing = computedStyle.letterSpacing;
  clone.style.textTransform = computedStyle.textTransform;
  clone.style.position = "fixed";
  clone.style.top = spineRect.top + "px";
  clone.style.left = spineRect.left + "px";
  clone.style.width = spineRect.width + "px";
  clone.style.height = slot.height + "px";
  clone.style.transform = "none";
  clone.style.margin = "0";
  clone.style.zIndex = "99";
  document.body.appendChild(clone);
  const slotRect = slot.getBoundingClientRect();
  const stripRect = strip.getBoundingClientRect();
  let dy = stripRect.left - slotRect.left;
  clone.style.transition = "all 1s ease";
  clone.style.top = slotRect.top + "px";
  clone.style.left = slotRect.left + "px";
  clone.style.width = dy + "px";
  await wait2(1e3);
  clone.style.transition = "all 0.25s ease-in";
  clone.style.transform = " scale(0.88)";
  clone.style.transformOrigin = "center top";
  await wait2(500);
  clone.style.transition = "all 0.15s ease-out";
  clone.style.transform = "scale(0.9)";
  flickerLed(romled);
  await wait2(500);
  await closeCap();
  cleapUpOverlay();
};
var flickerLed = async (led) => {
  led.classList.remove("off");
  await wait2(60);
  led.classList.add("off");
  await wait2(80);
  led.classList.remove("off");
  await wait2(150);
  led.classList.add("off");
  await wait2(120);
  led.classList.remove("off");
};
var cleapUpOverlay = () => {
  overlay.style.display = "none";
  middle.classList.remove("rotated");
  middle.classList.add("active");
};
var drawNav = () => {
  var c = document.getElementById("canvas-back");
  var ctx3 = c.getContext("2d");
  var PS = 16, COLS = 14, ROWS = 14;
  var T = { hi: "#ffc040", md: "#b05000", dk: "#3a1400" };
  var G = [0, 0, 0, 0, 0, 0, 0, "hi", "hi", "hi", "hi", "hi", "hi", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "md", 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "md", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "dk", "dk", "dk", "dk", "hi", 0];
  for (var i = 0; i < G.length; i++) {
    if (!G[i]) continue;
    var r = Math.floor(i / COLS), col = i % COLS;
    ctx3.fillStyle = T[G[i]];
    ctx3.fillRect(col * PS, r * PS, PS, PS);
  }
  var c = document.getElementById("canvas-next");
  var ctx3 = c.getContext("2d");
  var PS = 16, COLS = 14, ROWS = 14;
  var T = { hi: "#ffc040", md: "#b05000", dk: "#3a1400" };
  var G = [0, "hi", "hi", "hi", "hi", "hi", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "md", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, "md", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "dk", "dk", "dk", "dk", "dk", 0, 0, 0, 0, 0, 0, 0];
  for (var i = 0; i < G.length; i++) {
    if (!G[i]) continue;
    var r = Math.floor(i / COLS), col = i % COLS;
    ctx3.fillStyle = T[G[i]];
    ctx3.fillRect(col * PS, r * PS, PS, PS);
  }
};

// static/scripts/nes.js
var currentIndex = 3;
var middle2 = document.getElementById("middle");
var left = document.getElementById("left");
var right = document.getElementById("right");
var slot2 = document.getElementById("slot");
var overlay2 = document.getElementById("overlay");
var strip2 = document.getElementById("strip");
var cap2 = document.getElementById("cartridge");
var GamesArray = [];
var powerLed = document.getElementById("power");
var powerBtn = document.getElementById("start");
var pauseBtn = document.getElementById("pause");
var romLoaded = "";
var power = false;
var roms = [left, middle2, right];
document.addEventListener("DOMContentLoaded", async () => {
  await makeMockTiles();
  startRandomTipEngine();
  await setUpKnobs();
  await setUpButtons();
  await initGames();
  await wait(500);
  updateRom();
});
var BeginKnobs = async () => {
  knobSettings["sound"].setting = 4;
  knobSettings["angle"] = 180;
  await turnKnob("sound", 180);
  await turnKnob("speed", 180);
};
var setUpKnobs = async () => {
  const soundKnob = document.getElementById("sound");
  const soundIncrements = [0, 45, 90, 135, 180, 225, 270, 315, 360];
  soundKnob.addEventListener("click", async () => {
    if (romLoaded == "" || power == false) {
      await wiggleKnob("sound", 45);
      await turnKnob("sound", 0);
      return;
    }
    knobSettings["sound"].setting = knobSettings["sound"].setting + 1;
    if (knobSettings["sound"].setting >= soundIncrements.length - 1) {
      knobSettings["sound"].setting = 0;
    }
    let targetAngle = soundIncrements[knobSettings["sound"].setting];
    if (targetAngle == 0) {
      updateVolume(0);
    } else {
      updateVolume(targetAngle / 180);
    }
    await turnKnob("sound", targetAngle);
  });
  const speedKnob = document.getElementById("speed");
  const speedIncrements = [45, 90, 135, 180, 225, 270, 315, 360];
  speedKnob.addEventListener("click", async () => {
    if (romLoaded == "" || power == false) {
      await wiggleKnob("speed", 45);
      await turnKnob("speed", 0);
      return;
    }
    knobSettings["speed"].setting = knobSettings["speed"].setting + 1;
    if (knobSettings["speed"].setting >= speedIncrements.length - 1) {
      knobSettings["speed"].setting = 0;
    }
    let targetAngle = speedIncrements[knobSettings["speed"].setting];
    UpdateSpeed(1e3 * targetAngle / 180);
    await turnKnob("speed", targetAngle);
  });
};
var initGames = async () => {
  let response = await fetch("./games");
  if (response.ok) {
    let data = await response.json();
    data.forEach((element) => {
      GamesArray.push(element);
    });
  }
};
function updateRom() {
  let romSelection = GamesArray.slice(currentIndex - 1, currentIndex + 2);
  if (currentIndex == 0) {
    romSelection = [GamesArray.at(-1), GamesArray[0], GamesArray[1]];
  }
  if (currentIndex == GamesArray.length - 1) {
    romSelection = [GamesArray.at(-2), GamesArray.at(-1), GamesArray[0]];
  }
  if (romSelection.length == 3) {
    for (let i = 0; i < 3; i++) {
      const element = romSelection[i];
      const rom = roms[i];
      rom.id = element.ID;
      const rom_title = rom.querySelector("span");
      const img = rom.querySelector("img");
      img.src = "/rom_images/" + element.ID + ".webp";
      const spine = rom.querySelector(".rom-spine");
      spine.innerText = element.name;
      rom_title.innerText = element.name;
    }
  }
}
var setUpButtons = async () => {
  pauseBtn.addEventListener("click", async () => {
    if (pauseBtn.classList.contains("paused")) {
      pauseBtn.classList.remove("paused");
      pauseBtn.classList.add("unpaused");
      pauseBtn.innerHTML = `   <svg viewBox="0 0 24 24" fill="#5a3a1a" stroke="#5a3a1a" stroke-width="1.5" stroke-linecap="square" xmlns="http://www.w3.org/2000/svg">
                                     <rect x="4" y="3" width="4" height="18"/>
                                     <rect x="16" y="3" width="4" height="18"/></svg>`;
      await ResumeGame();
    } else {
      pauseBtn.classList.remove("unpaused");
      pauseBtn.classList.add("paused");
      pauseBtn.innerHTML = ` <svg viewBox="0 0 24 24" fill="#5a3a1a" stroke="#5a3a1a" stroke-width="1.5"
                                    stroke-linecap="square">
                                    <polygon points="6,3 20,12 6,21" />
                                </svg>`;
      await PauseGame();
    }
  });
  powerBtn.addEventListener("click", async () => {
    power = true;
    if (romLoaded != "") {
      begin(romLoaded);
      return;
    }
    powerLed.classList.remove("off");
    await openCap(cap2);
    await wait(200);
    overlay2.style.display = "flex";
    overlay2.style.transition = "all ease 1";
    overlay2.style.opacity = 1;
  });
  const nav_back = document.getElementById("nav-back");
  nav_back.addEventListener("click", async () => {
    move(-1);
    pushbtn("canvas-back");
  });
  const nav_front = document.getElementById("nav-front");
  nav_front.addEventListener("click", async () => {
    move(1);
    pushbtn("canvas-next");
  });
  middle2.addEventListener("click", async () => {
    middle2.classList.remove("active");
    await wait(50);
    let span = middle2.querySelector("span");
    span.style.display = "none";
    middle2.classList.add("rotated");
    romLoaded = middle2.id;
    await slotCart();
    if (power) {
      activateJoypad();
      await begin(romLoaded);
    }
  });
  cap2.addEventListener("click", async () => {
    await PauseGame();
    await openCap(cap2);
    await wait(200);
    overlay2.style.display = "flex";
    overlay2.style.transition = "all ease 1";
    overlay2.style.opacity = 1;
  });
};
async function begin(game) {
  initConsole();
  await loadRom(game);
  await BeginKnobs();
  alignCable(1);
  state.romRunning = true;
}
function move(direction) {
  let newIdx = currentIndex + direction;
  if (newIdx < 0) {
    currentIndex = GamesArray.length - 1;
    updateRom();
    return;
  }
  if (newIdx >= GamesArray.length) {
    currentIndex = newIdx - GamesArray.length;
    updateRom();
    return;
  }
  currentIndex = newIdx;
  updateRom();
}
