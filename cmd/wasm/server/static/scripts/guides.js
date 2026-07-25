import { marked } from "./marked.esm.js"

const guideWindow = document.querySelector(".guide-overlay")
const markdownContainer = document.querySelector(".markdown-body")

export const openGuides = async ()=>{ 
    await fetchAndFill("sample")
}

const fetchAndFill =  async (name)=>{
    const response = await fetch(`./docs/${name}.md`)
    const md = await response.text()

    console.log(marked.parse(md))
    markdownContainer.innerHTML = md
    
}