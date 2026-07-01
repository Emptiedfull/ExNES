(() => {
  // static/scripts/emuWorker.js
  var instance = null;
  importScripts("/wasm_exec.js");
  var init = async () => {
    const go = new Go();
    const result = await WebAssembly.instantiateStreaming(fetch("/nes.wasm"), go.importObject);
    go.run(result.instance);
    instance = result.instance;
    self.postMessage({ "type": "wasm" });
    reset();
  };
  init();
  var FBuf = new SharedArrayBuffer(256 * 240 * 4);
  var FBytes = new Uint8Array(FBuf);
  var SIZE = 8192;
  var audioBufS = new SharedArrayBuffer(SIZE * 4 + 4 + 4 + 4 + 4);
  var samples = new Float32Array(audioBufS, 0, SIZE);
  var control = new Int32Array(audioBufS, SIZE * 4, 3);
  var S_size = 1024;
  var S_buf = new Float32Array(S_size);
  self.onmessage = async ({ data }) => {
    switch (data.type) {
      case "init":
        startEmulator();
        initBuffer(new Uint8Array(S_buf.buffer));
        initInput(new Int32Array(data.inputBuf));
        initSpeed(new Uint8Array(data.speedBuf));
        self.postMessage({ type: "init", audioBufS, FBuf, SIZE, S_size });
        break;
      case "loadRom":
        console.log("loading rom");
        await loadRom(data.rom);
        pump();
        break;
      case "pump":
        pump();
        break;
      case "reset":
        console.log("Recieved reset request");
        reset();
        break;
      case "getsnap":
        console.log("getting snapshot list");
        let snapshots = await requestSnapshotList();
        self.postMessage({ type: "snaps", snaps: snapshots });
        break;
      case "loadSnapshot":
        loadSnapshot(data.index);
        break;
      case "input":
        update(data.action, data.pressed);
        break;
    }
  };
  var requestSnapshotList = async () => {
    const headerPTR = getSnapshotList();
    const header = new Uint32Array(instance.exports.memory.buffer, headerPTR, 2);
    const snapshotPTR = header[0];
    const len = header[1];
    const snapshotBYTES = new Uint8Array(instance.exports.memory.buffer, snapshotPTR, len);
    console.log(new TextDecoder().decode(snapshotBYTES));
    const snapshots = JSON.parse(new TextDecoder().decode(snapshotBYTES));
    return snapshots;
  };
  var pump = () => {
    while (true) {
      Atomics.wait(control, 2, 0);
      if (Atomics.load(control, 2) == 2) {
        return;
      }
      const wp = Atomics.load(control, 0);
      const rp = Atomics.load(control, 1);
      const free = SIZE - (wp - rp);
      const want = Math.min(free, S_size);
      if (want > 0) {
        drive(want);
        for (let i = 0; i < want; i++) {
          samples[(wp + i) % SIZE] = S_buf[i];
        }
        Atomics.store(control, 0, wp + want);
      }
      FBytes.set(new Uint8Array(frameBuffer.buffer));
      self.postMessage({ type: "frameUp" });
      Atomics.store(control, 2, 0);
      Atomics.notify(control, 2);
    }
  };
  var loadRom = async (game) => {
    const response = await fetch("/games/" + game + ".nes");
    if (!response.ok) {
      createModal("Unable to fetch Rom", response, true);
      return;
    }
    const buffer = await response.arrayBuffer();
    const uint8view = new Uint8Array(buffer);
    initRom(uint8view);
  };
})();
