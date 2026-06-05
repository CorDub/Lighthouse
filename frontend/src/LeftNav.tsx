import "./styles/LeftNav.css"
import AddButton from "./AddButton.tsx";

function LeftNav() {
  function addReport() {
    
  }

  return (
    <div class="leftNav">
      <AddButton clickFn={addReport}/>
    </div>
  )
}

export default LeftNav;