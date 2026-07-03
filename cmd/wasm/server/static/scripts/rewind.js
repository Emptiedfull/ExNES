import { wait } from "./joypad"
import { loadSnap } from "./driver"

const track = document.getElementById("items-container")
const preview_img = document.getElementById("preview-img")
const rewind_overlay = document.getElementById("rewind-overlay")

const timeline = document.getElementById("rewind-timeline")
const highligh_view = document.getElementById("rewind-highlight")


const back = document.getElementById("rewind-back")
const next = document.getElementById("rewind-next")
const start = document.getElementById("rewind-start")
const end = document.getElementById("rewind-end")

const load = document.getElementById("review-load")
const cart = document.getElementById("review-cart")

const saveCap = document.getElementById("rewind-cap")
const saveArea = document.getElementById("rewind-area")
const slot = document.getElementById("rewind-slot")

let activeID = 0

const clicker = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {

            clickSound()


            makeActive(entry.target)
            updateHighlight()

        }
    })
},
    {
        root: track,
        rootMargin: '0px -45% 0px -45%',
        threshold: 0.1,
    })

let ctx = new AudioContext()




export const createTilesFromSnapshots = async (snapshots) => {
    console.log("creating snapshot tiles")
    openCap()
    let lastTile = null
    rewind_overlay.style.display = "flex"
    setupRewindButtons(snapshots.length)
    setupTimeline(snapshots.length)
    for (let i = 0; i < snapshots.length; i++) {
        let snapshot = snapshots[i]
        let tile = document.createElement("div")
        tile.id = "tile-" + i
        clicker.observe(tile)

        tile.classList.add("tape-item")

        let img = document.createElement("canvas")

        img.height = snapshot.image.height
        img.width = snapshot.image.width

        tile.addEventListener("mousedown", () => {
            makeActive(tile)
        })

        img.getContext('2d').drawImage(snapshot.image, 0, 0)

        lastTile = tile
        tile.append(img)
        track.append(tile)
    }
    
    lastTile.scrollIntoView({
        behavior: "smooth",
        inline: "nearest",
        block: "center"
    })

}

export const makeMockTiles = async () => {
    console.log("mocking ur ass")
    

   

    // saveCap.style.setProperty("offset","-100%")

    let lastTile = null
    setupRewindButtons(100)
    setupTimeline(100)
    for (let i = 0; i < 100; i++) {
        let tile = document.createElement("div")
        tile.classList.add("tape-item")
        tile.id = "tile-" + i

        let img = createRandomPixelArt()
        tile.appendChild(img)

        tile.addEventListener("mousedown", () => {
            makeActive(tile)
        })

        clicker.observe(tile)
        lastTile = tile

        track.append(tile)

    }
    lastTile.scrollIntoView({
        behavior: "smooth",
        inline: "nearest",
        block: "center"
    })

}

const setupRewindButtons = (length) => {
    next.addEventListener("click", () => {
        let target = parseInt(activeID, 10) + 1
        if (target > length - 1) {
            target = length - 1
        }
        let targetEl = document.getElementById("tile-" + target)
        makeActive(document.getElementById("tile-" + target))
        scrollReelToValue(target, convertValToPerc().center)

    })

    back.addEventListener("click", () => {
        let target = parseInt(activeID, 10) - 1
        if (target < 0) {
            target = 0
        }

        let targetEl = document.getElementById("tile-" + target)
        makeActive(targetEl)
        scrollReelToValue(target, convertValToPerc().center)
    })

    end.addEventListener("click",async()=>{
        let target = length - 1
         let targetEl = document.getElementById("tile-" + target)
          scrollReelToValue(target, false)
          await wait(200)
        makeActive(targetEl)
       
    })

    start.addEventListener("click",async()=>{
        let target = 0
         let targetEl = document.getElementById("tile-" + target)
          scrollReelToValue(target, false)
          await wait(200)
        makeActive(targetEl)
       
    })

    load.addEventListener("click",async()=>{
        console.log("loading snapshot:",activeID)
       

        

        await startSaveLoad()
         await loadSnap(parseInt(activeID,10))

      
    
    })
} 


const startSaveLoad = async () => {
  let clone = cart.cloneNode(true)
  clone.classList.add("rewind-cart-clone")

  let originalRect = cart.getBoundingClientRect()
  let slotRect = slot.getBoundingClientRect()
  let areaRect = saveArea.getBoundingClientRect()
  let areaHeight = areaRect.top - areaRect.bottom
  let sourceCanvas = cart.querySelector("canvas")

  let width = originalRect.right - originalRect.left
  let targetWidth = slotRect.right - slotRect.left
  
  let scaleFactor = targetWidth / width;

  
  clone.style.position = "fixed"
  clone.style.left = originalRect.left + "px"
  clone.style.top = originalRect.top + "px"
  clone.style.width = width + "px"
  clone.style.transformOrigin = "top left"
  clone.style.transition = "transform 0.8s cubic-bezier(0.34, 1.56, 0.67 , 1)"
  clone.style.zIndex = 0

  clone.querySelector("canvas").getContext("2d").drawImage(sourceCanvas, 0, 0);

  document.body.appendChild(clone);

  await wait(10);
      rewind_overlay.style.display = "none"
  
  let dx = (slotRect.left ) - originalRect.left;
  let dy = (slotRect.top ) - originalRect.top;

  clone.style.transform = `translate(${dx}px, ${dy}px) scale(${scaleFactor})`;

  await wait(1000)

  let movedown = slotRect.bottom - slotRect.top
    clone.style.transition = "transform 0.8s ease"
  
   clone.style.transform = `translate(${dx}px, ${dy - areaHeight}px) scale(${scaleFactor})`;
    closeCap()

  
};

const openCap = ()=>{
    console.log(slot)
    let Rect = saveArea.getBoundingClientRect()
    let dy = Rect.top - Rect.bottom 
    console.log(dy)

    let capRect = saveCap.getBoundingClientRect()
    let capheight = capRect.bottom - capRect.top
    

    saveArea.style.transform = `translateY(${dy}px)`
    saveCap.style.transform = `translateY(${dy - capheight}px)`
}

const closeCap =()=>{
    saveArea.style.transform = `translateY(0)`
    saveCap.style.transform =  `translateY(0)`
}


let lastactive = null

const makeActive = (el) => {


    if (lastactive != null) {
        lastactive.classList.remove("tape-active")
    }

    let source = el.querySelector("canvas")

    preview_img.height = 240
    preview_img.width = 256

    preview_img.getContext("2d").drawImage(source, 0, 0)

    el.classList.add("tape-active")

    let index = el.id.slice(5)

    timeline.value = index
    activeID = index

    lastactive = el
}

const setupTimeline = (val) => {
    timeline.min = 0
    timeline.max = val - 1
}

timeline.addEventListener("input", (e) => {

    

    scrollReelToValue(timeline.value, convertValToPerc().center)
})

const updateHighlight = () => {

    let res = convertValToPerc()

    highligh_view.style.left = res.perc - 6 + "%"
    return res.center
}



const scrollReelToValue = (value, center) => {


    let elemtentID = "tile-" + Math.round(value)

    let tile = document.getElementById(elemtentID)

    tile.scrollIntoView({
        behavior: "instant",
        inline: center ? "center" : "nearest",
        block: "center"
    })
}

const convertValToPerc = () => {
    let max = timeline.max
    let min = timeline.min
    let val = timeline.value

    let perc = ((val - min) / (max - min)) * 100

    let center = true

    if (perc < 6) {
        perc = 6
        center = false
    }

    if (perc > 94) {
        perc = 94
        center = false
    }


    return { perc, center }
}

const getCenter = (perc) => {

}


const clickSound = () => {

    if (ctx == null) return
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()

    osc.connect(gain)
    gain.connect(ctx.destination)

    osc.frequency.value = 900
    gain.gain.setValueAtTime(0.4, ctx.currentTime)
    gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.02)

    osc.start()
    osc.stop(ctx.currentTime + 0.02)

}



export const StartRewindEngine = async () => {
    rewind_overlay.style.display = "flex"


    await createTiles()
}

function createRandomPixelArt(blockSize = 4) {
    const canvas = document.createElement('canvas');
    canvas.width = 256;
    canvas.height = 240;
    const ctx = canvas.getContext('2d');


    const palette = ['#000000', '#FCFCFC', '#F8F8F8', '#BCBCBC',
        '#7C7C7C', '#A40000', '#0000FC', '#00A800',
        '#F8B800', '#00FCFC', '#F800F8', '#585858'];

    for (let y = 0; y < canvas.height; y += blockSize) {
        for (let x = 0; x < canvas.width; x += blockSize) {
            ctx.fillStyle = palette[Math.floor(Math.random() * palette.length)];
            ctx.fillRect(x, y, blockSize, blockSize);
        }
    }

    return canvas;
}