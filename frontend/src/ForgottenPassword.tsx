import "./ForgottenPassword.css";
import { createSignal } from "solid-js";
import Navbar from "./Navbar";

function ForgottenPassword() {
  const [email, setEmail] = createSignal("")

  async function resetPassword(e: Event) {
    try {
      e.preventDefault()
      const response = await fetch("/api/checkPassword", {
        method: "POST",
        credentials: "include",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
          email: email()
        })
      })

      if (response.ok) {
        const data = await response.json()
        console.log('data', data)
      }
    } catch(error) {
      console.error(error)
    }
  }

  return (
    <div class="forgotten-password">
      <Navbar/>
      <form 
        class="confirm-email-form"
        onSubmit={(e) => resetPassword(e)}>
        <input
          class="form-input"
          type="text"
          placeholder="Enter email address"
          onChange={(e) => setEmail(e.target.value)}/>
        <button 
          class="green-button clickable"
          onClick={(e) => resetPassword(e)}>
            Submit
        </button>
      </form>
    </div>
  )
}

export default ForgottenPassword;