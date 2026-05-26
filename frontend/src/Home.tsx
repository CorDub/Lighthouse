import { createSignal, Show, For } from "solid-js"
import Navbar from "./Navbar.tsx"
import { BASE_URL } from "./helpers/config.ts";
import Text from "./Text.tsx";
import { useDefaults } from "./DefaultsContext.tsx";
import homeText from "./translations/Home.json";

function Home () {
  const [users, setUsers] = createSignal(null);
  const { defaults } = useDefaults();

  async function fetchUsers() {
    const response = await fetch(`${BASE_URL}/api/users`, {
      method:"GET",
      credentials: "include",
      headers: {"Content-Type":"application/json"},
    })

    if (response.ok) {
      const data = await response.json()
      setUsers(data.users)
    }
  }

  return (
    <div class="home">
      <Navbar />
      <p><Text
          value={homeText.welcome} 
          lang={defaults().lang}/>
      </p>
      <button 
        class='accept-button clickable'
        onClick={fetchUsers}>
        <Text
          value={homeText.fetchButton} 
          lang={defaults().lang}/>
      </button>
      <Show
        when={users()}
      > 
        <For each={users()}>{(user, _) =>
          <p>{user.id}</p>
        }</For>
        <button 
          class="accept-button clickable"
          onClick={() => setUsers(null)}>
          <Text
            value={homeText.clear} 
            lang={defaults().lang}/>
        </button>
      </Show>
    </div>
  )
}

export default Home
