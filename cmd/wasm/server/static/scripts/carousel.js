let currentIndex = 0;
const track = document.querySelector('.romCarousel')

function move(direction) {
  const items = document.querySelectorAll('.romItem')
  currentIndex += direction
  
 
  currentIndex = Math.max(0, Math.min(currentIndex, items.length - 1))

  
 
  const offset = -currentIndex * 440; 
  track.style.transform = `translateX(${offset}px)`
  console.log("moving",direction)
}