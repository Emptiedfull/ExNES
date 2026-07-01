import { wait } from "./joypad"

const track = document.getElementById("items-containers")
const preview_img = document.getElementById("preview-img")
const rewind_overlay = document.getElementById("rewind-overlay")

const clicker = new IntersectionObserver((entries)=>{
    entries.forEach(entry=>{
        if (entry.isIntersecting){
            
            clickSound()
        }
    })
},
{
    root: track,
    rootMargin: '0px -45% 0px -45%', 
    threshold:0.1,
})

let ctx = null

// document.addEventListener("keydown",async()=>{
//     if (ctx == null){
//         ctx = new AudioContext()

//         await StartRewindEngine()
//     }   
// })


export const createTiles = async ()=>{
    const cont = document.getElementById("items-container")

    let el = null
    let lasttile = null

    for (let i = 0;i < 100;i++){
        let tile = document.createElement("div")
        let img = new Image()
        img.src = "./rom_images/tets.png"
        tile.appendChild(img)
        tile.id = "tile-" + i
        tile.classList.add("tape-item")
        clicker.observe(tile)

        tile.addEventListener("mousedown",()=>{
                makeActive(tile)

                tile.classList.add("tape-active")
            
        })

        cont.appendChild(tile)

        

        lasttile = tile
        tile.classList.remo
    }

    
    // preview_img.src =  "./rom_images/tets.png"


    lasttile.scrollIntoView({
        behavior:"smooth",
        inline:"nearest",
        block:"center"
    })

}

let lastactive = null

const makeActive = (el)=>{

    if (lastactive != null){
        lastactive.classList.remove("tape-active")
    }

    console.log("scrolling to:",el)
    preview_img.src = el.querySelector("img").src
    el.scrollIntoView({
        behavior:"smooth",
        inline:"nearest",
        block:"center"
    })

    lastactive = el
}


const clickSound = ()=>{
    
    if (ctx == null) return
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()

    osc.connect(gain)
    gain.connect(ctx.destination)

    osc.frequency.value = 900
    gain.gain.setValueAtTime(0.4,ctx.currentTime)
    gain.gain.exponentialRampToValueAtTime(0.001,ctx.currentTime + 0.02)

    osc.start()
    osc.stop(ctx.currentTime + 0.02)
 
}


export const StartRewindEngine = async ()=>{
    rewind_overlay.style.display = "flex"

    await createTiles()
}