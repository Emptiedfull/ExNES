import { wait } from "./joypad"

const track = document.getElementById("items-container")
const preview_img = document.getElementById("preview-img")
const rewind_overlay = document.getElementById("rewind-overlay")

const timeline = document.getElementById("rewind-timeline")

const clicker = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {

            clickSound()
        }
    })
},
    {
        root: track,
        rootMargin: '0px -45% 0px -45%',
        threshold: 0.1,
    })

let ctx = null

// document.addEventListener("keydown",async()=>{
//     if (ctx == null){
//         ctx = new AudioContext()

//         await StartRewindEngine()
//     }   
// })

export const createTilesFromSnapshots = async (snapshots) => {
    rewind_overlay.style.display = "flex"
    console.log(snapshots)
    for (let i = 0; i < snapshots.length; i++) {
        let snapshot = snapshots[i]
        let tile = document.createElement("div")

        tile.classList.add("tape-item")

        let img = document.createElement("canvas")

        img.height = snapshot.image.height
        img.width = snapshot.image.width

        tile.addEventListener("mousedown", () => {
            makeActive(tile, snapshot.image)
        })

        img.getContext('2d').drawImage(snapshot.image, 0, 0)


        tile.append(img)
        track.append(tile)
    }

}

export const makeMockTiles = async () => {
    setupTimeline(150)
    for (let i = 0; i < 100; i++) {
        let tile = document.createElement("div")
        tile.classList.add("tape-item")
        tile.id = "tile-" + i

        let img = document.createElement("canvas")
        tile.appendChild(img)

        tile.addEventListener("mousedown", () => {
            makeActive(tile,null,100)
        })

        track.append(tile)

    }
}


let lastactive = null

const makeActive = (el, image,length) => {

    if (lastactive != null) {
        lastactive.classList.remove("tape-active")
    }


    if (image != null) {
        preview_img.height = image.height
        preview_img.width = image.width
        preview_img.getContext("2d").drawImage(image, 0, 0)
    }

    el.scrollIntoView({
        behavior: "smooth",
        inline: "nearest",
        block: "center"
    })

    el.classList.add("tape-active")

    let index = el.id.slice(5)
    
    timeline.value = index

    lastactive = el
}

const setupTimeline = (val)=>{
    timeline.min = 0 
    timeline.max = val
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