import "./styles/LeftNav.css"
import AddButton from "./AddButton.tsx";
import { BASE_URL } from "./helpers/config.ts";
import { useUser } from "./UserContext.tsx";

function LeftNav() {
  const { user } = useUser()

  async function addReport() {
    const response = await fetch(`${BASE_URL}/api/reports`, {
      method: "POST",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify({
        userId: user()?.id
      })
    })

    if (response.ok) {
      const data = await response.json()
      console.log("data", data)
    }
  }

  return (
    <div class="leftNav">
      <AddButton clickFn={addReport}/>
    </div>
  )
}

export default LeftNav;