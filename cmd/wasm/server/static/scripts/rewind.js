import { wait } from "./joypad"

export const createTiles = async ()=>{
    const cont = document.getElementById("items-container")

    let el = null

    for (let i = 0;i < 100;i++){
        let tile = document.createElement("div")
        let img = new Image()
        img.src = "./rom_images/tets.png"
        tile.appendChild(img)
        tile.classList.add("tape-item")
        cont.appendChild(tile)

        if (i == 99){
            el = tile
        }

    }

    el.scrollIntoView({
        behavior:"smooth",
        inline:"nearest",
        block:"center"
    })

    // reelScrollTo(cont,el)

    clickSound()

}


const clickSound = ()=>{
    console.log("hello")
    const ctx = new AudioContext()
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()

    osc.connect(gain)
    gain.connect(ctx.destination)

    osc.frequency.value = 500
    gain.gain.setValueAtTime(0.9,ctx.currentTime)
    gain.gain.exponentialRampToValueAtTime(0.001,ctx.currentTime + 0.03)

    osc.start()
    osc.stop(ctx.currentTime + 0.03)

}
