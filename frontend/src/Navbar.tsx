import "./Navbar.css"
import { Show } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { useUser } from "./UserContext";

function Navbar() {
  const navigate = useNavigate()
  const { user, setUser } = useUser();

  async function logout() {
    const response = await fetch("/api/logout", {
      method: "POST",
      credentials: "include",
      headers: {"Content-Type": "application/json"}
    })

    if (response.ok) {
      setUser(null)
      navigate("/")
    }
  }

  return (
    <div class="navbar">
      <div class="nav-logo">Lighthouse</div>
      <Show when={user()}>
        <div class="nav-logout">
          <button 
            class="clickable white-button"
            onClick={logout}>
              Logout
            </button>
        </div>
      </Show>
    </div>
  )
}

export default Navbar