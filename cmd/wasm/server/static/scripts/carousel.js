let currentIndex = 0

const romArray = [{"name":"mario","id":"mario","img":"x"},{"name":"zelda","id":"zelda","img":"x"},{"name":"donkey","id":"donkey","img":"x"},{"name":"balloon","id":"balloon","img":"x"},{"name":"contra","id":"contra","img":"x"},{"name":"sb3","id":"sb3","img":"x"}]

const middle = document.getElementById("middle")
const left = document.getElementById("left")
const right = document.getElementById("right")  

const roms = [left,middle,right]



function move(direction) {
    newIdx = currentIndex + direction 
    if (newIdx >= romArray.length || newIdx < 0){
        return
    }

     currentIndex = newIdx

    romSelection = romArray.slice(currentIndex-1,currentIndex+2)

    if (romSelection.length == 3){
        for (let i = 0; i < 3; i++) {
            const element = romSelection[i];

            console.log(roms[i],element)
            
        }
    }

}


(function(){
  var c=document.getElementById('canvas-back');
  console.log("working")
  var ctx=c.getContext('2d');
  var PS=16, COLS=14, ROWS=14;
  var T={hi:'#d8d4b0',md:'#8a8660',dk:'#2a2810'};
  var G=[0,0,0,0,0,0,0,'hi','hi','hi','hi','hi','hi',0,0,0,0,0,0,0,'hi','md','md','md','md','dk',0,0,0,0,0,0,0,'hi','md','md','md','md','dk',0,0,0,0,0,0,0,'hi','md','md','md','md','dk',0,0,0,0,0,0,0,'hi','md','md','md','md','dk',0,0,0,0,0,0,0,'hi','md','md','md','md','dk',0,0,0,0,0,0,0,'hi','md','md','md','md','md',0,0,0,0,0,0,0,0,'dk','md','md','md','md','md',0,0,0,0,0,0,0,0,0,'dk','md','md','md','md','hi',0,0,0,0,0,0,0,0,0,'dk','md','md','md','md','hi',0,0,0,0,0,0,0,0,0,'dk','md','md','md','md','hi',0,0,0,0,0,0,0,0,0,'dk','md','md','md','md','hi',0,0,0,0,0,0,0,0,0,'dk','md','md','md','md','hi',0,0,0,0,0,0,0,0,0,'dk','dk','dk','dk','dk','hi',0];
  for(var i=0;i<G.length;i++){
    if(!G[i])continue;
    var r=Math.floor(i/COLS), col=i%COLS;
    ctx.fillStyle=T[G[i]];
    ctx.fillRect(col*PS,r*PS,PS,PS);
  }
})();


(function(){
  var c=document.getElementById('canvas-next');
  var ctx=c.getContext('2d');
  var PS=16, COLS=14, ROWS=14;
  var T={hi:'#ffc040',md:'#b05000',dk:'#3a1400'};
  var G=[0,'hi','hi','hi','hi','hi','hi',0,0,0,0,0,0,0,0,0,'dk','md','md','md','md','hi',0,0,0,0,0,0,0,0,0,'dk','md','md','md','md','hi',0,0,0,0,0,0,0,0,0,'dk','md','md','md','md','hi',0,0,0,0,0,0,0,0,0,'dk','md','md','md','md','hi',0,0,0,0,0,0,0,0,0,'dk','md','md','md','md','hi',0,0,0,0,0,0,0,0,0,'md','md','md','md','md','hi',0,0,0,0,0,0,0,0,'md','md','md','md','md','dk',0,0,0,0,0,0,0,'hi','md','md','md','md','dk',0,0,0,0,0,0,0,'hi','md','md','md','md','dk',0,0,0,0,0,0,0,'hi','md','md','md','md','dk',0,0,0,0,0,0,0,'hi','md','md','md','md','dk',0,0,0,0,0,0,0,'hi','md','md','md','md','dk',0,0,0,0,0,0,0,'hi','dk','dk','dk','dk','dk',0,0,0,0,0,0,0];
  for(var i=0;i<G.length;i++){
    if(!G[i])continue;
    var r=Math.floor(i/COLS), col=i%COLS;
    ctx.fillStyle=T[G[i]];
    ctx.fillRect(col*PS,r*PS,PS,PS);
  }
})();


document.addEventListener("DOMContentLoaded",async()=>{
    console.log(middle)
    middle.addEventListener("click",()=>{
        console.log("slotting")
        console.log(middle.classList)
    })
});

const slot = async ()=>{
    console.log("slotting")

    console.log(middle.classList)
}

move(1)