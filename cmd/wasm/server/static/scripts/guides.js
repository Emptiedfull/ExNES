import { marked } from "./marked.esm.js"

const guideWindow = document.querySelector(".guide-overlay")
const markdownContainer = document.querySelector(".markdown-body")

const closeBtn = document.getElementById("guide-close")

const guidesCont = document.getElementById("guideList")

export const openGuides = async ()=>{ 

    await fetchAndFillGuides()
   
    
}

const fetchAndFillGuides = async ()=>{
    const response = await fetch(`./docs`)
    const guides = await response.json()

   

    guides.forEach(element => {

        const guide = document.createElement("div")

        guide.innerText = element
        guide.classList.add("guide-item")
        guide.id = element
        
        if (element == "Introduction"){
            guide.classList.add("active")
            fetchAndFill(element)
        }

        guide.addEventListener("click",()=>{
         
            fetchAndFill(element)
        })

        guidesCont.appendChild(guide)
    });

    console.log(guides)
}

const makeActive = (name) =>{

    for (const child of parent.children){
        console.log(child)
    }

    const item = document.getElementById(name)


}

const fetchAndFill =  async (name)=>{
    makeActive(element)
    const response = await fetch(`./docs/${name}.md`)
    const md = await response.text()

    console.log(marked.parse(md))
    markdownContainer.innerHTML = marked.parse(md)
    
}