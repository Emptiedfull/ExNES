export const SetUpModals = ()=>{

}

const createModal = ()=>{
    modalList = document.querySelector(".model-bar")

    const modal = document.createElement("div")
    modal.classList.add("modal","modal-corners")

    const title = document.createElement("h1")
    title.classList.add("info")

    title.textContent = "HELLO THIS IS GOD"

    modal.appendChild(title)

    const para = document.createElement("p")
    para.classList.add("modal-text")
    modal.appendChild(para)

    para.textContent = "SORRY THIS IS NOT GOD"

}