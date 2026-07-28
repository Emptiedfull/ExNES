import { marked } from "./marked.esm.js"

const guideWindow = document.querySelector(".guide-overlay")
const markdownContainer = document.querySelector(".markdown-body")

const closeBtn = document.getElementById("guide-close")

const guidesCont = document.getElementById("guideList")

const docHeaders = []

export const openGuides = async ()=>{ 
     guideWindow.style.display = "flex"
     initButtons()
    await fetchAndFillGuides()
   
    
}

const initButtons = ()=>{
    closeBtn.addEventListener("click",()=>{

        guideWindow.style.display = "none"
    })
}

export const openGuideTarget = ()=>{
    guideWindow.computedStyleMap.op
}

export const closeGuide = ()=>{
    guideWindow.style.display = "none"
}



const fetchAndFillGuides = async ()=>{
    const response = await fetch(`./docs`)
    const guides = await response.json()

   

    guides.forEach(element => {

        const guide = document.createElement("div")



        guide.innerText = element
        guide.classList.add("guide-item")
        guide.id = element

        docHeaders.push(element)
        
        if (element == "Introduction"){
            guide.classList.add("active")
            fetchAndFill(element)
        }

        guide.addEventListener("click",()=>{
         
            fetchAndFill(element)
        })

        guidesCont.appendChild(guide)
        console.log(guidesCont)
    });
    
}

const makeActive = (name) =>{

    for (const child of guidesCont.children){
        console.log(child)
        if (child.classList.contains("active")){
            child.classList.remove("active")
        }
    }

    const item = document.getElementById(name)
    if (item){
         item.classList.add("active")
    }
   


}

const fetchAndFill =  async (name)=>{
    makeActive(name)
    const response = await fetch(`./docs/${name}.md`)
    const md = await response.text()

   
    markdownContainer.innerHTML = marked.parse(md)
     hijackLinks(marked.parse(md))
    
   markdownContainer.scrollTo({ top: 0, behavior: "smooth" })
}

const hijackLinks = (content)=>{

    const links = document.querySelectorAll("a")
    const start = window.location.href.length

    links.forEach(element =>{
        const dest = element.href.slice(start)
       

        if (dest.includes("docs")){
            const doc = dest.slice(5,-3)
          

            if (docHeaders.includes(doc)){
               
                element.addEventListener("click",(e)=>{
                    e.preventDefault()
                    console.log("clciked me hehe")

                    fetchAndFill(doc)
                })
            }
          
           
        
        }

        
        
        
    })

}