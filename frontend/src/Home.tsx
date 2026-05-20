import { createSignal, Show, For } from "solid-js"
import Navbar from "./Navbar.tsx"
import { BASE_URL } from "./helpers/config.ts";

function Home () {
  const [users, setUsers] = createSignal(null);

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
      <p>Yes this is home</p>
      <button 
        class='accept-button clickable'
        onClick={fetchUsers}>
        Fetch users
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
          Clear
        </button>
      </Show>
    </div>
  )
}

export default Home
