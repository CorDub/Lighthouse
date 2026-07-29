import "./styles/LeftNav.css"
import AddButton from "./AddButton.tsx";
import { setDisplayed } from "./stores/agencyHomeDisplay.tsx";

function LeftNav() {
  return (
    <div class="leftNav">
      <AddButton clickFn={setDisplayed} arg={"createReport"}/>
    </div>
  )
}

export default LeftNav;