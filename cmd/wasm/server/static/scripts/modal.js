import { wait } from "./joypad"

export const SetUpModals = ()=>{
    

}

export const createModal = async (head,body,error = false,fleeting = true)=>{

    
    let modalList = document.querySelector(".moodle-bar")

    const modal = document.createElement("div")
    modal.classList.add("modal","modal-corners")
    
    const title = document.createElement("h1")

    if (error){
        title.classList.add("error")
    }else{
        title.classList.add("info")
    }
    
    title.textContent = head
    modal.appendChild(title)

    const para = document.createElement("p")
    para.classList.add("modal-text")
    para.textContent = body
    modal.appendChild(para)

    modalList.appendChild(modal)

    if (fleeting){
        setUpForFailure(modal)
    }
    
    activateModal(modal)

    return modal
}


export const activateModal = (modal)=>{
    modal.classList.add("active")
}

export const deactivateModal = async (modal) =>{
    modal.classList.remove("active")
    modal.classList.add("inactive")

    await wait(400)
    modal.remove()

}

const setUpForFailure = async (modal)=>{
    await wait(3500)

    if (modal.classList.contains("active")){
        await deactivateModal(modal)
    }else{
        console.log("bro alr failed")
    }
}