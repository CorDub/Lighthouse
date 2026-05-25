import "./styles/Navbar.css"
import { Show } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { useUser } from "./UserContext";
import { BASE_URL } from "./helpers/config.ts";
import LanguagePicker from "./LanguagePicker.tsx";

type NavbarProps = {
  lang: "en" | "es";
  setLang: (lang: "en" | "es") => void;
}

function Navbar(props: NavbarProps) {
  const navigate = useNavigate()
  const { user, setUser } = useUser();

  async function logout() {
    const response = await fetch(`${BASE_URL}/api/logout`, {
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
      <div class="nav-logo">
        <a 
          class="subdued-white-link"
          onClick={() => navigate("/")}>
            Lighthouse
        </a>
      </div>
      <div class="nav-right">
        <LanguagePicker lang={props.lang} setLang={props.setLang}/>
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
    </div>
  )
}

export default Navbar